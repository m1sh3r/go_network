package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithBasicAuth(t *testing.T) {
	h := withBasicAuth(getUsers)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()
	h(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusUnauthorized)
	}

	req = httptest.NewRequest(http.MethodGet, "/users", nil)
	req.SetBasicAuth("admin", "secret")
	rr = httptest.NewRecorder()
	h(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusOK)
	}
}
