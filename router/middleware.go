package router

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/72sevenzy2/json-parser/helpers"
)

type Middleware func(HandlerFunc) HandlerFunc // the middleware type (takes in the current handler and returns a new one)

// custom responseWriter type to capture status code and request byte size.
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// override WriteHeader func()
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code // saving status code in struct
	rw.ResponseWriter.WriteHeader(code)
}

// overriding Write to capture bytes
func (rw *responseWriter) Write(b []byte) (int, error) {
	v, err := rw.ResponseWriter.Write(b)
	rw.size += v // tracking the size in bytes
	return v, err
}

// plain structs to work with default values.
type bodySize struct {
	size uint32
}

// // object pooling for request body logs
// var buff = sync.Pool{
// 	New: func() any {
// 		return new(bytes.Buffer)
// 	},
// }

func Logger(confSize uint32) Middleware { // returns the middleware type
	return func(hf HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *Request) {
			start := time.Now() // setting the current time (before the request has ended)
			fmt.Printf("Request has started with URL: %s, and method: %s, and in time: %s\n", r.URL, r.Method, start)

			// buffer for comparison in limited writer Write().
			buf := bytes.Buffer{}

			// setting default value first for request body size

			var opt bodySize
			// only set if confSize was set to 0 (will indicate to user in docs):
			if confSize == 0 {
				opt = bodySize{
					size: 1024, // 1kb
				}
			} else {
				// apply custom size
				opt = bodySize{
					size: confSize,
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

			fmt.Println("request body data: (with data size of:)", opt.size)
			fmt.Println(string(body))

			// redacting sensitive header before printing
			header := r.Header.Clone()
			header.Del("Authorization")

			fmt.Println("Request headers:", header)
		}
	}
}

// bearer auth middleware (this includes having a bearer token which will then be compared to the authkey )

func BearerAuth(AuthKey string) Middleware {
	return func(hf HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *Request) {
			if len(AuthKey) <= 1 { // check if authkey has less than 1 character
				helpers.Failed(w, http.StatusInternalServerError, errors.New("please include a stronger AuthKey."))
				return
			}

			authLab := r.Header.Get("Authorization") // grabbing the token

			var token string
			if v := strings.Contains(AuthKey, "Bearer "); v {
				token = strings.TrimPrefix(authLab, "Bearer ") // removing the "bearer " part of the token to then compare it to the authkey

				if token != AuthKey {
					helpers.Failed(w, http.StatusForbidden, errors.New("invalid token."))
					return
				}

				hf(w, r) // next handler
			}

			// continuing if "Bearer " doesnt include in the authkey.
			if AuthKey != authLab { // check if the authkey is matching
				helpers.Failed(w, http.StatusForbidden, errors.New("invalid token.")) // if not then throw a failed json response
				return                                                               // exit the request
			}

			hf(w, r) // continue to next handler
		}
	}
}

// basic auth middleware (this auth includes having a user and password inorder to access the endpoint)

func BasicAuth(user, password string) Middleware { // implements the middleware type which returns a handler
	return func(hf HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *Request) {

			authUser, authPassword, ok := r.BasicAuth() // extracting the user and password and if it exists (ok) from the r.BasicAuth() func, which is a built in method in go to do so, instead of manually parsing it ourselves.

			if !ok || authUser != user || authPassword != password { // run the necessary logic
				helpers.Failed(w, http.StatusForbidden, errors.New("invalid credentials."))
				return
			}

			hf(w, r) // continue to next handler
		}
	}
}

// timeout middleware

func Timeout(seconds int) Middleware {
	return func(hf HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *Request) {
			ctx, cancel := context.WithTimeout(r.Context(), time.Duration(seconds)*time.Second) // initialising timeout (in seconds)

			defer cancel() // cancelling at the end of the func (current handler)

			hf(w, &Request{
				contextReq: r.WithContext(ctx), // pass in with request context (creates shallow copy of the original request and runs it with context)
			}) // ServeHTTP(w, and "r" with the context 'ctx')
		}
	}
}

// recoverer middleware (for preventing server crashes)

func Recoverer() Middleware {
	return func(hf HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *Request) {
			defer func() { // catches any crashses and recovers the request, while printing the err in return.
				if err := recover(); err != nil {
					fmt.Println("caught: ", err)
					helpers.Failed(w, http.StatusInternalServerError, fmt.Errorf("server error: %s", err))
				}
			}()

			hf(w, r) // next handler
		}
	}

}

// middleware chaining.

// func to apply the middlewares
func (r *Router) ApplyMiddlewares(h HandlerFunc) HandlerFunc {
	for i := len(r.Middlewares) - 1; i >= 0; i-- {
		h = r.Middlewares[i](h)
	}

	return h
}

// // Use func to use the middewares (also appending it to the Middlewares type in router struct
// func (r *Group) Use(s Middleware) {
// 	r.mws = append(r.mws, s)
// }
