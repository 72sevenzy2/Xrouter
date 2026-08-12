package router

import (
	"net/http"

	"github.com/72sevenzy2/http-router/core"
)

// handler methods as alternative to Handle(...).

// Router methods
func (r *Router) Get(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	r.Handle(http.MethodGet, path, handler, mws...)
}

func (r *Router) Post(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	r.Handle(http.MethodPost, path, handler, mws...)
}

func (r *Router) Put(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	r.Handle(http.MethodPut, path, handler, mws...)
}

func (r *Router) Del(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	r.Handle(http.MethodDelete, path, handler, mws...)
}
// connect, head, options, patch, tree
func (r *Router) Connect(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	r.Handle(http.MethodConnect, path, handler, mws...)
}

func (r *Router) Head(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	r.Handle(http.MethodHead, path, handler, mws...)
}

func (r *Router) Options(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	r.Handle(http.MethodOptions, path, handler, mws...)
}

func (r *Router) Patch(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	r.Handle(http.MethodPatch, path, handler, mws...)
}

func (r *Router) Trace(path string, handler core.HandlerFunc, mws ...core.Middleware) {
	r.Handle(http.MethodTrace, path, handler, mws...)
}
