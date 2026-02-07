package main

import (
	"fmt"
	"net/http"
)

type userInfo struct {
	Password string
	Role     string
}

var users = map[string]userInfo{
	"admin":  {Password: "secret", Role: "admin"},
	"editor": {Password: "edit123", Role: "editor"},
	"viewer": {Password: "view123", Role: "viewer"},
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	info, exists := users[user]
	if !exists || info.Password != pass {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	switch info.Role {
	case "admin":
		fmt.Fprintln(w, "Добро пожаловать, администратор. Полный доступ.")
	case "editor":
		fmt.Fprintln(w, "Здравствуйте, редактор. Доступ на редактирование.")
	case "viewer":
		fmt.Fprintln(w, "Здравствуйте, наблюдатель. Только чтение.")
	default:
		fmt.Fprintf(w, "Здравствуйте, %s.\n", user)
	}
}

func main() {
	http.HandleFunc("/admin", protectedHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Не удалось запустить сервер: %v\n", err)
	}
}
