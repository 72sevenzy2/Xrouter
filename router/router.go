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
	"net/http"
	"strings"

	"github.com/72sevenzy2/json-parser/helpers"
)

// Route struct to loop over dynamic routes
type Route struct {
	Method  string
	Parts   []string
	handler http.HandlerFunc
}

// Router struct to hold all static/dynamic routes
type Router struct {
	StaticRoutes  map[string]map[string]http.HandlerFunc
	DynamicRoutes []Route
	Middlewares   []Middleware // storing our middlewares here (type is our Middleware function type)
}

// util functions

// func to check whether route is dynamic or static
func isDynamic(path string) bool {
	return strings.Contains(path, ":") // strings.Contains() returns a boolean
}

// func to split route path
func splitPath(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
	/*
				Splits the route path as so:

				example route: /users/42

				strings.Trim(path, "/") splits any slashes from start to end from the route,
				and strings.Split() splits the route from slashes ANYWHERE, so the route afterwards this block of code would look like:
		 		{"users", "42"}, as strings.Split() returns as string array
	*/
}

func NewRouter() *Router {
	// contructing the router upon the func being called
	return &Router{
		StaticRoutes: make(map[string]map[string]http.HandlerFunc), // initialising the map of map)
	} // which is just: "PATH": "...": "METHOD": ... (method can be either get, post, put, etc)
}

// adding routes, and assigning the method of the route aswell as the url to the handler which then is executed in the ServeHTTP func
func (r *Router) Handle(method, path string, handler http.HandlerFunc, mws ...Middleware) {

	// applying route specific middleware in reverse order
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler) // add handler to mw
	}

	// check if route is dynamic
	if isDynamic(path) {
		r.DynamicRoutes = append(r.DynamicRoutes, Route{
			Method:  method,
			Parts:   splitPath(path),
			handler: handler,
		})
		return
	}

	// otherwise use static route logic
	if r.StaticRoutes[path] == nil { // check if route already exists before creating
		r.StaticRoutes[path] = make(map[string]http.HandlerFunc)
	}

	r.StaticRoutes[path][method] = handler // assign both static route path and method to handler


	// keeping commented code for future reference
	// if r.StaticRoutes[path] == nil { // checking if route endpoint itself doesnt exist before creating the path and assigning it to another map (map[string]http.HandlerFunc) which will be for the method map
	// 	r.StaticRoutes[path] = make(map[string]http.HandlerFunc) // assign the path to the method type (GET, POST, PUT etc)
	// }

	// for i := len(mws) -1; i >= 0; i-- { // applying the middleware in reverse order
	// 	handler = mws[i](handler) // assign the handler to the middleware (the middleware returns a new handler after performing its programmed task.)
	// }  // assigning it to the index of [i] which is each middleware that gets passed in as the mws parameter, so it can be 1 or many.

	// r.StaticRoutes[path][method] = handler // assign method to the handler (handler type is http.handlerFunc)
	// // we're basically taking the path which will be something like "/hi": and the method name, or its type we can call it
	// // for example:       "/hi":
	// //                       "GET": "and some handler here, (in this case, it will be the http handlerfunc we used)"
}

// core routing logic for my router
func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// attempt static routes (lower time complexity compared to dynamicRoutes which requires looping over routes (O(n)))
	if methods, ok := s.StaticRoutes[r.URL.Path]; ok { // validate if route path exists
		if handler, ok := methods[r.Method]; ok { // also check if method for route path is appropriate
			finalHandler := s.ApplyMiddlewares(handler) // apply middlewares if included
			finalHandler(w, r) // run handler
			return
		}

		// otherwise return err if method is invalid
		helpers.Failed(w, http.StatusMethodNotAllowed, "method not allowed")
		return 
	}

	// attempt dynamic routes (higher time complexity than staticRoutes (which are O(1)))

	// parts := splitPath(r.URL.Path)

	for _, route := range s.DynamicRoutes {
		if route.Method != r.Method {
			continue // skipping the current iteration
		}


	}

	// refactored code ^ (keeping commented code for future reference)

	// method, ok := s.StaticRoutes[r.URL.Path] // search for if the url path exists

	// if !ok {
	// 	http.NotFound(w, r) // check if path exists on the request being sent.
	// 	return
	// }

	// handler, ok := method[r.Method] // assign the handler to the method type

	// if !ok { // and check if the method of the path is valid.
	// 	helpers.Failed(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	// 	return
	// }

	// finalHandler := s.ApplyMiddlewares(handler)
	// finalHandler(w, r)
}
