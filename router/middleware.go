package router

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/72sevenzy2/json-parser/helpers"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc // the middleware type (takes in the current handler and returns a new one)

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

// object pooling for request body logs
var buff = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func Logger(confSize uint32) Middleware { // returns the middleware type
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now() // setting the current time (before the request has ended)
			fmt.Printf("Request has started with URL: %s, and method: %s, and in time: %s\n", r.URL, r.Method, start)

			buf := buff.Get().(*bytes.Buffer) // reusable buffer.
			buf.Reset()                       // reset buffer before use.
			defer buff.Put(buf)               // add buf to the buffer pool

			// setting default value first for request body size
			opt := bodySize{
				size: 1024, // 1kb
			}

			// apply custom size
			if confSize != 0 {
				opt = bodySize{
					size: confSize,
				}
			}

			// limit size
			lm := &LimitedBuffer{
				buf:   buf,
				limit: opt.size + 1,
			}

			r.Body = io.NopCloser(io.TeeReader(r.Body, lm)) // using io.NopCloser as io.TeeReader does not implement io.ReadCloser.
			// io.TeeReader allows the current handler to read the request body data, whilst also allowing copying.

			rw := &responseWriter{ // default status code and custom response writer initialisation
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			hf(rw, r) // calling the next function to continue to the next handler
			// by calling hf() before printing, we give time to the io Readers above to read the request body data.

			// compressing body if over 1 kb
			body := buf.Bytes()
			if uint32(len(body)) > opt.size {
				body = body[:opt.size]; // truncated
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
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if len([]byte(AuthKey)) == 1 { // check if authkey has more than 1 character (converting to []byte gives a string each char 1 byte.)
				helpers.Failed(w, http.StatusInternalServerError, "please include a stronger AuthKey.")
				return
			}

			authLab := r.Header.Get("Authorization") // grabbing the token

			token := strings.TrimPrefix(authLab, "Bearer ") // removing the "bearer " part of the token to then compare it to the authkey

			if token == authLab || token != AuthKey { // check if the authkey is matching
				helpers.Failed(w, http.StatusForbidden, "Invalid Token") // if not then throw a failed json response
				return                                                   // exit the request
			}

			hf(w, r) // continue to next handler
		}
	}
}

// basic auth middleware (this auth includes having a user and password inorder to access the endpoint)

func BasicAuth(user, password string) Middleware { // implements the middleware type which returns a handler
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			authUser, authPassword, ok := r.BasicAuth() // extracting the user and password and if it exists (ok) from the r.BasicAuth() func, which is a built in method in go to do so, instead of manually parsing it ourselves.

			if !ok || authUser != user || authPassword != password { // run the necessary logic
				helpers.Failed(w, http.StatusForbidden, "invalid credentials.")
				return
			}

			hf(w, r) // continue to next handler
		}
	}
}

// timeout middleware

func Timeout(seconds int) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), time.Duration(seconds)*time.Second) // initialising timeout (in seconds)

			defer cancel() // cancelling at the end of the func (current handler)

			hf(w, r.WithContext(ctx)) // ServeHTTP(w, and "r" with the context 'ctx')
		}
	}
}

// recoverer middleware (for preventing server crashes)

func Recoverer() Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() { // catches any crashses and recovers the request, while printing the err in return.
				if err := recover(); err != nil {
					fmt.Println("caught: ", err)
					helpers.Failed(w, http.StatusInternalServerError, fmt.Sprintf("server error: %s", err))
				}
			}()

			hf(w, r) // next handler
		}
	}

}

// func to apply the middlewares
func (r *Router) ApplyMiddlewares(h http.HandlerFunc) http.HandlerFunc {
	for i := len(r.Middlewares) - 1; i >= 0; i-- {
		h = r.Middlewares[i](h)
	}

	return h
}

// Use func to use the middewares (also appending it to the Middlewares type in router struct
func (r *Router) Use(s Middleware) {
	r.Middlewares = append(r.Middlewares, s)
}
