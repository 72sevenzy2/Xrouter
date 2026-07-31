package router

import (
	"github.com/72sevenzy2/http-router/core"
)

// middleware chaining.
func (r *Router) ApplyMiddlewares(h core.HandlerFunc) core.HandlerFunc {
	for i := len(r.Middlewares) - 1; i >= 0; i-- {
		h = r.Middlewares[i](h)
	}

	return h
}
