package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSessionLoginProtectedLogout(t *testing.T) {
	loginReq := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"user","password":"user123"}`))
	loginRR := httptest.NewRecorder()
	loginHandler(loginRR, loginReq)

	if loginRR.Code != http.StatusOK {
		t.Fatalf("login code got %d, want %d", loginRR.Code, http.StatusOK)
	}

	cookies := loginRR.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("session cookie is missing")
	}

	protectedReq := httptest.NewRequest(http.MethodGet, "/protected", nil)
	for _, c := range cookies {
		protectedReq.AddCookie(c)
	}
	protectedRR := httptest.NewRecorder()
	protectedHandler(protectedRR, protectedReq)

	if protectedRR.Code != http.StatusOK {
		t.Fatalf("protected code got %d, want %d", protectedRR.Code, http.StatusOK)
	}
	if !strings.Contains(protectedRR.Body.String(), "user") {
		t.Errorf("unexpected protected body: %q", protectedRR.Body.String())
	}

	logoutReq := httptest.NewRequest(http.MethodGet, "/logout", nil)
	for _, c := range cookies {
		logoutReq.AddCookie(c)
	}
	logoutRR := httptest.NewRecorder()
	logoutHandler(logoutRR, logoutReq)

	if logoutRR.Code != http.StatusOK {
		t.Fatalf("logout code got %d, want %d", logoutRR.Code, http.StatusOK)
	}
}
