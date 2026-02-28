package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func makeToken(t *testing.T, username string, role string) string {
	t.Helper()
	claims := Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}
	return tokenString
}

func TestLoginHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"admin","password":"secret"}`))
	rr := httptest.NewRecorder()

	loginHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusOK)
	}

	var resp loginResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("token is empty")
	}
}

func TestProtectedHandlerByRole(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected string
	}{
		{"admin", "admin", "Панель администратора"},
		{"user", "user", "Зона пользователя"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := makeToken(t, tt.name, tt.role)
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+token)
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
