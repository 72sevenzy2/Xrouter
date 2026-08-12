/*
MIT License


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
	"slices"
	"strings"

	"github.com/72sevenzy2/http-router/core"
	"github.com/72sevenzy2/json-parser/helpers"
)

// Router struct to hold all static/dynamic routes
type Router struct {
	*core.Router
}

// type aliases core/types.go
type Request = core.Request

// Use func to use the middewares (also appending it to the Middlewares type in router struct
func (r *Router) Use(s core.Middleware) { // global
	r.Middlewares = append(r.Middlewares, s)
}

// initialise a new router.
func NewRouter() *Router {
	// contructing the router upon the func being called
	return &Router{
		Router: &core.Router{
			StaticRoutes:  make(map[string]map[string]core.HandlerFunc), // initialising the map of map)
			DynamicRoutes: make(map[string][]core.Route),
		},
	} // which is just: "PATH": "...": "METHOD": ... (method can be either get, post, put, etc)
}

// initial handler, (mainly for handlers.go)
func (r *Router) Handle(method string, path string, handler core.HandlerFunc, mws ...core.Middleware) {

	if len(mws) > 0 {
		// applying route specific middleware in reverse order
		for i := range slices.Backward(mws) {
			handler = mws[i](handler) // add handler to mw
		}
	}
	// check if route is dynamic
	if isDynamic(path) {
		r.DynamicRoutes[method] = append(r.DynamicRoutes[method], core.Route{
			Method:  method,
			Parts:   splitPath(path),
			Handler: handler,
		})
		return
	}

	// otherwise use static route logic
	if r.StaticRoutes[path] == nil { // check if route already exists before creating
		r.StaticRoutes[path] = make(map[string]core.HandlerFunc)
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
			finalHandler(w, &core.Request{
				Request:    r,
				ContextReq: r,
			}) // run handler
			return
		}

		// otherwise return err if method is invalid
		helpers.Failed(w)
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
		final(w, &core.Request{ // store params
			Request: r,

			Params: params,
		})
		return
	}

	helpers.Failed(w)

}
