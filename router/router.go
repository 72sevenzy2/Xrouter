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
	"errors"
	"net/http"
	"strings"

	"github.com/72sevenzy2/json-parser/helpers"
)

// custom request struct to hold routing essentials,
type Request struct {
	*http.Request
	contextReq *http.Request     // for timeout mw (WithContext() usage)
	params     map[string]string // holding route parameters
}

// custom handler type
type HandlerFunc func(http.ResponseWriter, *Request)

// Route struct to loop over dynamic routes
type Route struct {
	Handler HandlerFunc
	Parts   []string
	Method  string
}

// Router struct to hold all static/dynamic routes
type Router struct {
	DynamicRoutes map[string][]Route // split dynamic routes by methods to reduce lookup time
	Middlewares   []Middleware       // storing our middlewares here (type is our Middleware function type)
	StaticRoutes  map[string]map[string]HandlerFunc
}

// group type (for route grouping)
type Group struct {
	router *Router
	prefix string
	mws    []Middleware
}

// group method (adding a parent route)
func (r *Router) Group(p string) *Group {
	return &Group{
		router: r,
		prefix: p,
	}
}

// Use func to use the middewares (also appending it to the Middlewares type in router struct
func (r *Router) Use(s Middleware) { // global
	r.Middlewares = append(r.Middlewares, s)
}

// Group based middleware
func (g *Group) Use(s Middleware) {
	g.mws = append(g.mws, s)
}

// Handler func for grouped routes.
func (g *Group) Handle(method, path string, handler HandlerFunc, mws ...Middleware) {
	newPath := g.prefix + path // path included with parent route

	mw := append([]Middleware{}, g.mws...) // appending empty Middleware slice, and storing g.mws  (group based middleware)
	mw = append(mw, mws...) // appending route specific middleware

	g.router.Handle(method, newPath, handler, mw...)
}

// initialise a new router.
func NewRouter() *Router {
	// contructing the router upon the func being called
	return &Router{
		StaticRoutes:  make(map[string]map[string]HandlerFunc), // initialising the map of map)
		DynamicRoutes: make(map[string][]Route),
	} // which is just: "PATH": "...": "METHOD": ... (method can be either get, post, put, etc)
}

// adding routes, and assigning the method of the route aswell as the url to the handler which then is executed in the ServeHTTP func
func (r *Router) Handle(method string, path string, handler HandlerFunc, mws ...Middleware) {

	if len(mws) > 0 {
		// applying route specific middleware in reverse order
		for i := len(mws) - 1; i >= 0; i-- {
			handler = mws[i](handler) // add handler to mw
		}
	}

	// check if route is dynamic
	if isDynamic(path) {
		r.DynamicRoutes[method] = append(r.DynamicRoutes[method], Route{
			Method:  method,
			Parts:   splitPath(path),
			Handler: handler,
		})
		return
	}

	// otherwise use static route logic
	if r.StaticRoutes[path] == nil { // check if route already exists before creating
		r.StaticRoutes[path] = make(map[string]HandlerFunc)
	}

	r.StaticRoutes[path][method] = handler // assign both static route path and method to handler
}

// routing logic
func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// making all routes normalised without a / at the end, but with root /.
	// for example. Input: "/users//" output: "/users", input: "users/1" output: "/users/1"
	path := strings.TrimRight(r.URL.Path, "/")
	if path == "" {
		path = "/"
	}

	// attempt static routes (lower time complexity compared to dynamicRoutes which requires looping over routes (O(n)))
	if methods, ok := s.StaticRoutes[path]; ok { // validate if route path exists
		if handler, ok := methods[r.Method]; ok { // also check if method for route path is appropriate
			finalHandler := s.ApplyMiddlewares(handler) // apply middlewares if included
			finalHandler(w, &Request{
				contextReq: r,
			}) // run handler
			return
		}

		// otherwise return err if method is invalid
		helpers.Failed(w, http.StatusMethodNotAllowed, errors.New("method not allowed."))
		return
	}

	// attempt dynamic routes (higher time complexity than staticRoutes (which are O(1)))

	parts := splitPath(path)

	for _, route := range s.DynamicRoutes[r.Method] { // loop over dynamic routes which are grouped by methods

		ok, params := match(route.Parts, parts)
		if !ok {
			continue // skip current iteration
		}

		final := s.ApplyMiddlewares(route.Handler)
		final(w, &Request{ // store params
			Request: r,
			params:  params,
		})
		return
	}

	helpers.Failed(w, http.StatusNotFound, errors.New("Page not found."))

}
