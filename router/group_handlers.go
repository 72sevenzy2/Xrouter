package router

import (
	"net/http"

	"github.com/72sevenzy2/http-router/core"
)

// group handlers.

func (g *Group) Get(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	g.Handle(http.MethodGet, path, handler, mws...)
}

func (g *Group) Post(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	g.Handle(http.MethodPost, path, handler, mws...)
}
