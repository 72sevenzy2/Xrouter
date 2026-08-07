package router

import (
	"slices"

	mw "github.com/72sevenzy2/Xrouter-middlewares"
	"github.com/72sevenzy2/http-router/core"
)

// Limiter holds the NewLimiter() method in which is for initialising the struct, and RateLimiter() method in which is the middleware.
type Limiter = mw.Limiter 

// middleware chaining.
func (r *Router) ApplyMiddlewares(h core.HandlerFunc) core.HandlerFunc {
	for i := range slices.Backward(r.Middlewares) {
		h = r.Middlewares[i](h)
	}

	return h
}

func Logger(conf uint32) core.Middleware {
	return mw.Logger(conf)
}

func BasicAuth(user, pass string) core.Middleware {
	return mw.BasicAuth(user, pass)
}

func BearerAuth(token string) core.Middleware {
	return mw.BearerAuth(token)
}

func Recoverer() core.Middleware {
	return mw.Recoverer()
}

func Timeout(N int) core.Middleware {
	return mw.Timeout(N)
}
