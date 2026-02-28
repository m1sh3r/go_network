package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func resetSwaggerUsersState() {
	mu.Lock()
	defer mu.Unlock()
	users = map[int]User{
		1: {ID: 1, Name: "Alexey"},
		2: {ID: 2, Name: "Daria"},
	}
	nextID = 3
}

func TestDeleteUser(t *testing.T) {
	resetSwaggerUsersState()

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
