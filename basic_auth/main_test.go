package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProtectedHandlerUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rr := httptest.NewRecorder()

	protectedHandler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestProtectedHandlerRoles(t *testing.T) {
	tests := []struct {
		name     string
		user     string
		pass     string
		expected string
	}{
		{"admin", "admin", "secret", "администратор"},
		{"editor", "editor", "edit123", "редактор"},
		{"viewer", "viewer", "view123", "наблюдатель"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/admin", nil)
			req.SetBasicAuth(tt.user, tt.pass)
			rr := httptest.NewRecorder()

			protectedHandler(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("got %d, want %d", rr.Code, http.StatusOK)
			}
			if !strings.Contains(rr.Body.String(), tt.expected) {
				t.Errorf("response %q does not contain %q", rr.Body.String(), tt.expected)
			}
		})
	}
}
