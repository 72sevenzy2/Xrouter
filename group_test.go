package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/72sevenzy2/http-router/router"
	"github.com/72sevenzy2/http-router/core"
)

func TestStandardGrouping(t *testing.T) {
	b := router.NewRouter()

	api := b.Group("/test")

	api.Handle(http.MethodGet, "/test2", func(w http.ResponseWriter, r *core.Request) {
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

	api.Handle(http.MethodGet, "/testr", func(w http.ResponseWriter, r *core.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test/r1/r2/testr", nil)
	rr := httptest.NewRecorder()

	b.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Log(req.URL.Path)
	}
}

// nested grouping test

func TestNestedGroups(t *testing.T) {
	b := router.NewRouter()

	api := b.Group("/parent")

	v1 := api.Group("/child1")
	v1.Use(router.Logger(0))

	v2 := v1.Group("/child2")
	v2.Use(router.Recoverer())

	v2.Handle(http.MethodGet, "/testchild", func(w http.ResponseWriter, r *core.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRequest(http.MethodGet, "/parent/child1/child2/testchild", nil)
	req := httptest.NewRecorder()

	b.ServeHTTP(req, rr)

	if req.Code != http.StatusOK { // didnt cross endpoint
		t.Fatalf("failed on path: %s", rr.URL.Path)
	}
}
