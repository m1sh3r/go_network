package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type userInfo struct {
	Password string
	Role     string
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

var users = map[string]userInfo{
	"admin": {Password: "secret", Role: "admin"},
	"user":  {Password: "user123", Role: "user"},
}

const secretKey = "my-secret-key"

func main() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Некорректный JSON", http.StatusBadRequest)
			return
		}

		info, exists := users[req.Username]
		if !exists || info.Password != req.Password {
			http.Error(w, "Неверные учётные данные", http.StatusUnauthorized)
			return
		}

		claims := Claims{
			Username: req.Username,
			Role:     info.Role,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			http.Error(w, "Не удалось подписать токен", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(loginResponse{Token: tokenString}); err != nil {
			http.Error(w, "Не удалось сформировать ответ", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		// Получаем токен из заголовка Authorization
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Токен не передан", http.StatusUnauthorized)
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(tokenString, bearerPrefix) {
			http.Error(w, "Некорректный заголовок Authorization", http.StatusUnauthorized)
			return
		}
		tokenString = strings.TrimPrefix(tokenString, bearerPrefix)

		// Парсим токен
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неподдерживаемый алгоритм подписи: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			http.Error(w, "Токен не прошёл проверку", http.StatusUnauthorized)
			return
		}

		// Проверяем валидность токена
		if !token.Valid {
			http.Error(w, "Токен недействителен", http.StatusUnauthorized)
			return
		}

		switch claims.Role {
		case "admin":
			fmt.Fprintln(w, "Панель администратора")
		case "user":
			fmt.Fprintln(w, "Зона пользователя")
		default:
			http.Error(w, "Неизвестная роль", http.StatusForbidden)
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Не удалось запустить сервер: %v\n", err)
	}
}
