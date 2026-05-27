/*
MIT License

Copyright (c) 2026 72sevenzy2

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package router

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/72sevenzy2/json-parser/helpers"
)

// custom type for context label (for ServeHTTP() func)
type contextKey string

const paramsKey contextKey = "params"

// Route struct to loop over dynamic routes
type Route struct {
	Method  string
	Parts   []string
	Handler http.HandlerFunc
}

// Router struct to hold all static/dynamic routes
type Router struct {
	StaticRoutes  map[string]map[string]http.HandlerFunc
	DynamicRoutes []Route
	Middlewares   []Middleware // storing our middlewares here (type is our Middleware function type)
}

func NewRouter() *Router {
	// contructing the router upon the func being called
	return &Router{
		StaticRoutes: make(map[string]map[string]http.HandlerFunc), // initialising the map of map)
	} // which is just: "PATH": "...": "METHOD": ... (method can be either get, post, put, etc)
}

// adding routes, and assigning the method of the route aswell as the url to the handler which then is executed in the ServeHTTP func
func (r *Router) Handle(method string, path string, handler http.HandlerFunc, mws ...Middleware) {

	// applying route specific middleware in reverse order
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler) // add handler to mw
	}

	// check if route is dynamic
	if isDynamic(path) {
		r.DynamicRoutes = append(r.DynamicRoutes, Route{
			Method:  method,
			Parts:   splitPath(path),
			Handler: handler,
		})
		return
	}

	// otherwise use static route logic
	if r.StaticRoutes[path] == nil { // check if route already exists before creating
		r.StaticRoutes[path] = make(map[string]http.HandlerFunc)
	}

	r.StaticRoutes[path][method] = handler // assign both static route path and method to handler
}

// Grab parameters stored in context. Pass in an http.Request, and key (string) to be retrieved.
func Param(r *http.Request, key string) (string, bool) {
	params, ok := r.Context().Value(paramsKey).(map[string]string) // type asserting it to map[string]string, and getting the whole map.
	if !ok {
		return "", false
	}
	return params[key], true // retrieve a particular field
}

// core routing logic for my router
func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// making all routes normalised without a / at the end, but with root /.
	// for example. Input: "/users//" output: "/users", input: "users/1" output: "/users/1"
	path := strings.TrimRight(r.URL.Path, "/");
	if path == "" {
		path = "/" 
	}

	// attempt static routes (lower time complexity compared to dynamicRoutes which requires looping over routes (O(n)))
	if methods, ok := s.StaticRoutes[r.URL.Path]; ok { // validate if route path exists
		if handler, ok := methods[r.Method]; ok { // also check if method for route path is appropriate
			finalHandler := s.ApplyMiddlewares(handler) // apply middlewares if included
			finalHandler(w, r)                          // run handler
			return
		}

		// otherwise return err if method is invalid
		helpers.Failed(w, http.StatusMethodNotAllowed, errors.New("method not allowed."))
		return
	}

	// attempt dynamic routes (higher time complexity than staticRoutes (which are O(1)))

	parts := splitPath(r.URL.Path)

	for _, route := range s.DynamicRoutes {
		if route.Method != r.Method {
			continue // skipping the current iteration
		}

		ok, params := match(route.Parts, parts)
		if !ok {
			continue // skip current iteration
		}

		// apply and execute with context
		ctx := context.WithValue(r.Context(), paramsKey, params)
		final := s.ApplyMiddlewares(route.Handler)
		final(w, r.WithContext(ctx))
		return
	}

	helpers.Failed(w, http.StatusNotFound, errors.New("Page not found."))

}
