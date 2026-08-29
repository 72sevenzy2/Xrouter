package core

// this package will hold shared types from the router and the middlewares, to avoid import cycles upon instillation.

import (
	"net/http"
)

// Param is to store route parameters.
type Param struct {
	Key   string
	Value string
}
type Params []*Param
type Request struct { // shared
	*http.Request
	ContextReq *http.Request     // for timeout mw (WithContext() usage)
	Params     Params // holding route parameters
}

type HandlerFunc func(http.ResponseWriter, *Request) // shared
type Middleware func(HandlerFunc) HandlerFunc        // the middleware type (takes in the current handler and returns a new one)

// Route struct to loop over dynamic routes
type Route struct {
	Handler HandlerFunc
	Parts   []string
	Method  string
}

// Router struct to hold all static/dynamic routes
type Router struct { // shared
	DynamicRoutes map[string][]Route // split dynamic routes by methods to reduce lookup time
	Middlewares   []Middleware       // storing our middlewares here (type is our Middleware function type)
	StaticRoutes  map[string]map[string]HandlerFunc
}
