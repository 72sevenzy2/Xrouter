package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/72sevenzy2/http-router/router"
)

func TestRouter(t *testing.T) {
	b := router.NewRouter()

	b.Handle(http.MethodGet, "/test/", func(w http.ResponseWriter, r *router.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	rr := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ServeHTTP(rec, rr)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
