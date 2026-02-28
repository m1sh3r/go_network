package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func resetUsersState() {
	mu.Lock()
	defer mu.Unlock()
	users = map[int]User{
		1: {ID: 1, Name: "Alexey"},
		2: {ID: 2, Name: "Daria"},
	}
	nextID = 3
}

func TestParseUserID(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantID  int
		wantErr bool
	}{
		{"ok", "/users/5", 5, false},
		{"empty", "/users/", 0, true},
		{"bad", "/users/abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := parseUserID(tt.path)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if id != tt.wantID {
				t.Errorf("got %d, want %d", id, tt.wantID)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	resetUsersState()

	req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(`{"name":"Updated"}`))
	rr := httptest.NewRecorder()

	updateUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusOK)
	}

	var user User
	if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if user.Name != "Updated" {
		t.Errorf("got %q, want %q", user.Name, "Updated")
	}
}

func TestDeleteUser(t *testing.T) {
	resetUsersState()

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	rr := httptest.NewRecorder()

	deleteUser(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusNoContent)
	}

	mu.RLock()
	_, exists := users[1]
	mu.RUnlock()
	if exists {
		t.Error("user was not deleted")
	}
}
