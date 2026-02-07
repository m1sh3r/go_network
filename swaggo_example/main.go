package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	_ "akozadaev/swag_openAPI/docs"

	httpSwagger "github.com/swaggo/http-swagger"
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

// @Summary Get all users
// @Tags Users
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func getUsers(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Не удалось сформировать ответ", http.StatusInternalServerError)
	}
}

// @Summary Get user by ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 400 {string} string "Некорректный идентификатор"
// @Failure 404 {string} string "Пользователь не найден"
// @Router /users/{id} [get]
func getUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
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
		http.Error(w, "Не удалось сформировать ответ", http.StatusInternalServerError)
	}
}

// @Summary Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Success 201 {object} User
// @Failure 400 {string} string "Некорректный JSON"
// @Router /users [post]
func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	mu.Lock()
	newUser.ID = nextID
	users[nextID] = newUser
	nextID++
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newUser); err != nil {
		http.Error(w, "Не удалось сформировать ответ", http.StatusInternalServerError)
	}
}

// @Summary Delete user by ID
// @Tags Users
// @Param id path int true "User ID"
// @Success 204 {string} string "Deleted"
// @Failure 400 {string} string "Некорректный идентификатор"
// @Failure 404 {string} string "Пользователь не найден"
// @Router /users/{id} [delete]
func deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
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

// @title Users API
// @version 1.0
// @description A simple REST API for managing users
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getUsers(w, r)
		case "POST":
			createUser(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getUser(w, r)
		} else if r.Method == "DELETE" {
			deleteUser(w, r)
		} else {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("REST API запущен: http://localhost:8080")
	fmt.Println("Swagger UI доступен: http://localhost:8080/swagger/index.html")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Не удалось запустить сервер: %v\n", err)
	}
}
