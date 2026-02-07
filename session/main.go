package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = map[string]string{
	"admin": "secret",
	"user":  "user123",
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Используйте /login, /protected, /logout")
	})

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

		password, ok := users[req.Username]
		if !ok || password != req.Password {
			http.Error(w, "Неверные учётные данные", http.StatusUnauthorized)
			return
		}

		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, "Не удалось получить сессию", http.StatusInternalServerError)
			return
		}

		session.Values["username"] = req.Username
		if err := session.Save(r, w); err != nil {
			http.Error(w, "Не удалось сохранить сессию", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Вход выполнен")
	})

	http.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, "Не удалось получить сессию", http.StatusInternalServerError)
			return
		}

		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			http.Error(w, "Не авторизован", http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, "Здравствуйте, %s!", username)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, "Не удалось получить сессию", http.StatusInternalServerError)
			return
		}

		session.Options.MaxAge = -1
		if err := session.Save(r, w); err != nil {
			http.Error(w, "Не удалось очистить сессию", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Выход выполнен")
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Не удалось запустить сервер: %v\n", err)
	}
}
