package router

import (
	"github.com/72sevenzy2/http-router/core"
	mw "github.com/72sevenzy2/Xrouter-middlewares"
)

// todo: add middlewares as variables wrapped from middleware/...

// middleware chaining.
func (r *Router) ApplyMiddlewares(h core.HandlerFunc) core.HandlerFunc {
	for i := len(r.Middlewares) - 1; i >= 0; i-- {
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
