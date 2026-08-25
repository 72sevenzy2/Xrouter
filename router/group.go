package router

import (
	"net/http"
	"strings"

	"github.com/72sevenzy2/http-router/core"
)

// group type (for route grouping)
type Group struct {
	router *Router
	prefix string
	mws    []core.Middleware
}

// the router Group() method works such that when it is registered, child routes can also be registered using the parent route as it would be of type *Group, in which the child route would also inherit the parent routes middlewares.
// parent group method
func (r *Router) Group(p string, nests ...string) *Group {
	var updPath string
	updPath = p
	// make sure "/" is included in str
	if string(p[0]) != "/" {
		updPath = "/" + p
	}
	if len(nests) > 0 {
		for i := range nests {
			// Join() normalises each nests before assigning to updPath.
			updPath = Join(updPath, nests[i])
		}
	}

	return &Group{
		router: r,
		prefix: updPath,
	}
}

// Group based middleware
func (g *Group) Use(s core.Middleware) {
	g.mws = append(g.mws, s)
}

// group method for child routes
func (g *Group) Group(prefix string) *Group {
	cmws := append([]core.Middleware{}, g.mws...) // copy previous route nodes mw collection (each childing having their own mw slice to do so)

	return &Group{
		router: g.router,
		prefix: Join(g.prefix, prefix),
		mws:    cmws,
	}
}

// Handler func for grouped routes (for method-specific handlers in group_handlers.go)
func (g *Group) Handle(method, path string, handler core.HandlerFunc, mws ...core.Middleware) {
	if !IsValidHTTPMethod(method) {
		method = http.MethodGet // default to GET if invalid method.
	}
	normalisedpath := strings.TrimRight(path, "/")
	newPath := Join(g.prefix, normalisedpath) // path included as child route with parent route.

	mw := append([]core.Middleware{}, g.mws...) // appending empty Middleware slice, and storing g.mws  (group based middleware)
	mw = append(mw, mws...)                     // appending route specific middleware

	g.router.Handle(method, newPath, handler, mw...)
}
