package router

import (
	"net/http"
	"slices"

	"github.com/72sevenzy2/http-router/core"
)

// handler methods as alternative to Handle(...).

func (r *Router) Get(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	if len(mws) > 0 {
		for v := range slices.Backward(mws) {
			handler = mws[v](handler)
		}
	}

	if isDynamic(path) {
		r.DynamicRoutes[http.MethodGet] = append(r.DynamicRoutes[http.MethodGet], core.Route{
			Handler: handler,
			Method: http.MethodGet,
			Parts: splitPath(path),
		})
		return
	}

	// otherwise its static
	if r.StaticRoutes[path] == nil { // preventing duplicated routes.
		r.StaticRoutes[path] = make(map[string]core.HandlerFunc)
	}

	r.StaticRoutes[path][http.MethodGet] = handler
}

func (r *Router) Post(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	if len(mws) > 0 {
		for v := range slices.Backward(mws) {
			handler = mws[v](handler)
		}
	}

	if isDynamic(path) {
		r.DynamicRoutes[http.MethodPost] = append(r.DynamicRoutes[http.MethodPost], core.Route{
			Method: http.MethodPost,
			Handler: handler,
			Parts: splitPath(path),
		})
	}
	
	if r.StaticRoutes[path] == nil {
		r.StaticRoutes[path] = make(map[string]core.HandlerFunc)
	}
	r.StaticRoutes[path][http.MethodPost] = handler
}
