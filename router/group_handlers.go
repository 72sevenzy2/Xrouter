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

func (g *Group) Put(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	g.Handle(http.MethodPut, path, handler, mws...)
}

func (g *Group) Del(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	g.Handle(http.MethodDelete, path, handler, mws...)
}

func (g *Group) Trace(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	g.Handle(http.MethodTrace, path, handler, mws...)
}

func (g *Group) Connect(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	g.Handle(http.MethodConnect, path, handler, mws...)
}

func (g *Group) Head(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	g.Handle(http.MethodHead, path, handler, mws...)
}

func (g *Group) Patch(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	g.Handle(http.MethodPatch, path, handler, mws...)
}

func (g *Group) Options(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	g.Handle(http.MethodOptions, path, handler, mws...)
}
