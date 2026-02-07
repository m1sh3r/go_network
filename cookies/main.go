package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		const cookieName = "visit-count"
		count := 0

		// Чтение куки
		if cookie, err := r.Cookie(cookieName); err == nil {
			value := make(map[string]int)
			if err := cookieHandler.Decode(cookieName, cookie.Value, &value); err == nil {
				count = value["count"]
			}
		}

		count++

		// Устанавливаем куку с обновлённым счётчиком
		value := map[string]int{"count": count}
		encoded, err := cookieHandler.Encode(cookieName, value)
		if err != nil {
			http.Error(w, "Не удалось закодировать куку", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  cookieName,
			Value: encoded,
			Path:  "/",
		})

		fmt.Fprintf(w, "Вы посетили страницу %d раз", count)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Не удалось запустить сервер: %v\n", err)
	}
}
