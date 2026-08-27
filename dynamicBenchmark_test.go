package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/72sevenzy2/http-router/router"
)

func BenchmarkDynamicRoutes(t *testing.B) {
	r := router.NewRouter()

	r.Get("route/path-50/:id", func(w http.ResponseWriter, r *router.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/route/path-50/72", nil)
	rr := httptest.NewRecorder()
	t.ResetTimer()

	for t.Loop() {
		r.ServeHTTP(rr, req)
	}
}

// comparing to go's stdlib mux
func BenchmarkStdlibDynamic(t *testing.B) {
	r := http.NewServeMux()

	r.HandleFunc("route/path-50/:id", DummyHandler)

	req := httptest.NewRequest(http.MethodGet, "/route/path-50/72", nil)
	rr :=  httptest.NewRecorder()

	t.ResetTimer()

	for t.Loop() {
		r.ServeHTTP(rr, req)
	}
}

// shared with staticBenchmark_test.go
func DummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
