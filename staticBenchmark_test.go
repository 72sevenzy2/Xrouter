package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/72sevenzy2/http-router/router"
)

func BenchmarkStaticRoutes(t *testing.B) {
	r := router.NewRouter()

	r.Get("test/path/hi", func(w http.ResponseWriter, r *router.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test/path/hi", nil)
	rr := httptest.NewRecorder()

	t.ResetTimer()

	for t.Loop() {
		r.ServeHTTP(rr, req)
	}
}

func BenchmarkStdlibStatic(t *testing.B) {
	r := http.NewServeMux()

	r.HandleFunc("test/path/hi", DummyHandler)
	req := httptest.NewRequest(http.MethodGet, "/test/path/hi", nil)
	rr := httptest.NewRecorder()

	t.ResetTimer()
	for t.Loop() {
		r.ServeHTTP(rr, req)
	}
}
