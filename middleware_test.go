package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/72sevenzy2/http-router/router"
)

// auth testing
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

// test func to check if logger mw calls next middleware:
func TestLoggerNext(t *testing.T) {
	called := false

	next := func (w http.ResponseWriter, r *http.Request)  {
		called = true
		w.WriteHeader(http.StatusOK)
	}

	handler := router.Logger(1024)(next)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("test"))
	rr := httptest.NewRecorder()

	handler(rr, req)

	if !called {
		fmt.Println("logger did not call next().")
	}

	if rr.Code != http.StatusOK {
		t.Fatalf("failed with status %d", rr.Code)
	}

}

// test to make sure logger preserves data (body)
func TestLoggerBody(t *testing.T) {
	next := func (w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		if string(body) != "testC" {
			t.Fatalf("expected body %q, received %q", "testC", string(body))
		}
	}

	handler := router.Logger(1024)(next)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("testC"))
	rr := httptest.NewRecorder()

	handler(rr, req)
}