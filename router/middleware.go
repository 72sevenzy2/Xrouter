package router

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/72sevenzy2/json-parser/helpers"
	"github.com/72sevenzy2/http-router/core"
)


// custom responseWriter type to capture status code and request byte size.
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// override WriteHeader with custom response writer struct
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code // saving status code in struct
	rw.ResponseWriter.WriteHeader(code)
}

// overriding Write to capture byte size
func (rw *responseWriter) Write(b []byte) (int, error) {
	v, err := rw.ResponseWriter.Write(b)
	rw.size += v // tracking the size in bytes
	return v, err
}

// plain structs to work with default values.
type bodySize struct {
	size uint32
}

func Logger(confSize uint32) core.Middleware { // returns the middleware type
	return func(hf core.HandlerFunc) core.HandlerFunc {
		return func(w http.ResponseWriter, r *core.Request) {
			start := time.Now() // setting the current time (before the request has ended)
			fmt.Printf("Request has started with URL: %s, and method: %s, and in time: %s\n", r.URL, r.Method, start)

			// buffer for comparison in limited writer Write().
			buf := bytes.Buffer{}

			// setting default value first for request body size

			var opt bodySize
			// only set if confSize was set to 0 (will indicate to user in docs):
			opt = bodySize{
				size: 1024, // default
			}
			
			if confSize != 0 {
				opt = bodySize{
					size: confSize, // custom size
				}	
			}

			// limit size
			lm := &LimitedBuffer{
				buf:   &buf,
				limit: opt.size + 1,
			}

			r.Body = io.NopCloser(io.TeeReader(r.Body, lm)) // using io.NopCloser as io.TeeReader does not implement io.ReadCloser.
			// io.TeeReader allows the current handler to read the request body data, whilst also allowing copying.

			rw := &responseWriter{ // default status code and custom response writer initialisation
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			hf(rw, r)
			// by calling hf() before printing, we give time to the io Readers above to read the request body data.

			// truncating if over 1 kb
			body := buf.Bytes()
			if uint32(len(body)) > opt.size {
				body = body[:opt.size] // truncated
				fmt.Println("body has been truncated.")
			}

			endTime := time.Since(start) // after the request has ended, in which we will print below
			fmt.Printf("request has ended: %s, with status code %d ||| and with response body size (in bytes): %d", endTime, rw.status, rw.size)

			fmt.Println("\nrequest body data: (with data size of:)", opt.size)
			fmt.Println(string(body))

			// redacting sensitive header before printing
			header := r.Header.Clone()
			header.Del("Authorization")

			fmt.Println("Request headers:", header)
		}
	}
}

// bearer auth middleware (this includes having a bearer token which will then be compared to the authkey )

func BearerAuth(AuthKey string) core.Middleware {
	return func(hf core.HandlerFunc) core.HandlerFunc {
		return func(w http.ResponseWriter, r *core.Request) {
			if len(AuthKey) <= 1 { // check if authkey has less than 1 character
				helpers.Failed(w)
				return
			}

			authLab := r.Header.Get("Authorization") // grabbing the token

			var token string
			if v := strings.Contains(AuthKey, "Bearer "); v {
				token = strings.TrimPrefix(authLab, "Bearer ") // removing the "bearer " part of the token to then compare it to the authkey

				if token != AuthKey {
					helpers.Failed(w)
					return
				}

				hf(w, r) // next handler
			}

			// continuing if "Bearer " doesnt include in the authkey.
			if AuthKey != authLab { // check if the authkey is matching
				helpers.Failed(w) // if not then throw a failed json response
				return  // exit the request
			}

			hf(w, r) // continue to next handler
		}
	}
}

// basic auth middleware (this auth includes having a user and password inorder to access the endpoint)

func BasicAuth(user, password string) core.Middleware { // implements the middleware type which returns a handler
	return func(hf core.HandlerFunc) core.HandlerFunc {
		return func(w http.ResponseWriter, r *core.Request) {

			authUser, authPassword, ok := r.BasicAuth() // extracting the user and password and if it exists (ok) from the r.BasicAuth() func, which is a built in method in go to do so, instead of manually parsing it ourselves.

			if !ok || authUser != user || authPassword != password { // run the necessary logic
				helpers.Failed(w)
				return
			}

			hf(w, r) // continue to next handler
		}
	}
}

// timeout middleware

func Timeout(seconds int) core.Middleware {
	return func(hf core.HandlerFunc) core.HandlerFunc {
		return func(w http.ResponseWriter, r *core.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), time.Duration(seconds)*time.Second) // initialising timeout (in seconds)

			defer cancel() // cancelling at the end of the func (current handler)

			// shallow copy of original request, (preserving other Request{} fields)
			req := *r
			req.Request = r.WithContext(ctx)

			hf(w, &req)
		}
	}
}

// recoverer middleware (for preventing server crashes)

func Recoverer() core.Middleware {
	return func(hf core.HandlerFunc) core.HandlerFunc {
		return func(w http.ResponseWriter, r *core.Request) {
			defer func() { // catches any crashses and recovers the request, while printing the err in return.
				if err := recover(); err != nil {
					fmt.Println("caught: ", err)
					helpers.Failed(w)
				}
			}()

			hf(w, r) // next handler
		}
	}

}

// middleware chaining.

// func to apply the middlewares
func (r *Router) ApplyMiddlewares(h core.HandlerFunc) core.HandlerFunc {
	for i := len(r.Middlewares) - 1; i >= 0; i-- {
		h = r.Middlewares[i](h)
	}

	return h
}
