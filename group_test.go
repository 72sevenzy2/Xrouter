package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/72sevenzy2/http-router/router"
)

func TestStandardGrouping(t *testing.T) {
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

// test with inline nested routes

func TestInlineNests(t *testing.T) {
	b := router.NewRouter()

	api := b.Group("/test", "/r1", "/2")

	api.Handle(http.MethodGet, "/testr", func(w http.ResponseWriter, r *router.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test/r1/r2/testr", nil)
	rr := httptest.NewRecorder()

	b.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Log(req.URL.Path)
	}
}