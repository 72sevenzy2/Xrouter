package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/72sevenzy2/http-router/router"
)

func TestRouteGrouping(t *testing.T) {
	b := router.NewRouter()

	api := b.Group("/test")

	api.Handle(http.MethodGet, "/test2", func(w http.ResponseWriter, r *router.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test/test2", nil)

	rr := httptest.NewRecorder()

	b.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with path %s", rr.Code, req.URL.Path)
	}
}