package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/72sevenzy2/http-router/router"
)

func TestBasicAuth(t *testing.T) {
	b := router.NewRouter()

	// apply auth middleware
	b.Use(router.BasicAuth("user1", "pass1"))

	b.Handle(http.MethodGet, "/foo1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/foo1", nil)
	req.SetBasicAuth("user1", "pass1")

	rr := httptest.NewRecorder()

	b.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("failed with status %d", rr.Code)
	}
	fmt.Println("successful.")
}

func TestBearerAuth(t *testing.T) {
	b := router.NewRouter()

	b.Use(router.BearerAuth("bearerauth123"))

	b.Handle(http.MethodGet, "/foo2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/foo2", nil)
	// set auth header and key for BearerAuth()

	req.Header.Set("Authorization", "bearerauth123")

	rr := httptest.NewRecorder()

	b.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("failed auth with status %d", rr.Code)
	}
	fmt.Println("successful.")
}