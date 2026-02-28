package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVisitHandlerCounter(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	visitHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusOK)
	}
	if !strings.Contains(rr.Body.String(), "1") {
		t.Errorf("unexpected body: %q", rr.Body.String())
	}

	resp := rr.Result()
	cookies := resp.Cookies()
	if len(cookies) == 0 {
		t.Fatal("cookie is not set")
	}

	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.AddCookie(cookies[0])
	rr2 := httptest.NewRecorder()
	visitHandler(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", rr2.Code, http.StatusOK)
	}
	if !strings.Contains(rr2.Body.String(), "2") {
		t.Errorf("unexpected body: %q", rr2.Body.String())
	}
}
