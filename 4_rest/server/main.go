package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	users = map[int]User{
		1: {ID: 1, Name: "Alexey"},
		2: {ID: 2, Name: "Daria"},
	}
	nextID = 3
	mu     sync.RWMutex
)

var authUsers = map[string]string{
	"admin":  "secret",
	"reader": "read123",
}

func withBasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		expected, exists := authUsers[user]
		if !ok || !exists || expected != pass {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Не авторизован", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Не удалось закодировать список пользователей", http.StatusInternalServerError)
	}
}

func parseUserID(path string) (int, error) {
	idStr := strings.TrimPrefix(path, "/users/")
	if idStr == "" {
		return 0, errors.New("идентификатор не указан")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("некорректный идентификатор")
	}
	return id, nil
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseUserID(r.URL.Path)
	if err != nil {
		http.Error(w, "Некорректный идентификатор", http.StatusBadRequest)
		return
	}

	mu.RLock()
	user, exists := users[id]
	mu.RUnlock()

	if !exists {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Не удалось закодировать пользователя", http.StatusInternalServerError)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(newUser.Name) == "" {
		http.Error(w, "Поле name обязательно", http.StatusBadRequest)
		return
	}

	mu.Lock()
	newUser.ID = nextID
	users[nextID] = newUser
	nextID++
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newUser); err != nil {
		http.Error(w, "Не удалось закодировать пользователя", http.StatusInternalServerError)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseUserID(r.URL.Path)
	if err != nil {
		http.Error(w, "Некорректный идентификатор", http.StatusBadRequest)
		return
	}

	var update User
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(update.Name) == "" {
		http.Error(w, "Поле name обязательно", http.StatusBadRequest)
		return
	}

	mu.Lock()
	user, exists := users[id]
	if !exists {
		mu.Unlock()
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	user.Name = update.Name
	users[id] = user
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Не удалось закодировать пользователя", http.StatusInternalServerError)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseUserID(r.URL.Path)
	if err != nil {
		http.Error(w, "Некорректный идентификатор", http.StatusBadRequest)
		return
	}

	mu.Lock()
	if _, exists := users[id]; !exists {
		mu.Unlock()
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	delete(users, id)
	mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	http.HandleFunc("/users", withBasicAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getUsers(w, r)
		case "POST":
			createUser(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/users/", withBasicAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getUser(w, r)
		case "PUT":
			updateUser(w, r)
		case "DELETE":
			deleteUser(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}))

	fmt.Println("REST API запущен: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Не удалось запустить HTTP сервер: %v\n", err)
	}
}
