package router

import (
	"context"
	"slices"
	"time"

	mw "github.com/72sevenzy2/Xrouter-middlewares"
	"github.com/72sevenzy2/http-router/core"
)

// NewLimiter is a method of the Limiter struct, in which also holds the RateLimiter which is the middleware.
func NewLimiter(limit int, ref int) *mw.Limiter {
	return mw.NewLimiter(limit, ref)
}

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

func Canceller() (core.Middleware, context.CancelFunc) {
	return mw.Canceller()
}

func Recoverer() core.Middleware {
	return mw.Recoverer()
}

func Timeout(N time.Duration) core.Middleware {
	return mw.Timeout(N)
}
