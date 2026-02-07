package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type statusResponse struct {
	Status string `json:"status"`
}

type timeResponse struct {
	Time string `json:"time"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Привет, мир!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statusResponse{Status: "ok"}); err != nil {
		http.Error(w, "Не удалось сформировать ответ", http.StatusInternalServerError)
	}
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(timeResponse{Time: time.Now().Format(time.RFC3339)}); err != nil {
		http.Error(w, "Не удалось сформировать ответ", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/time", timeHandler)
	fmt.Println("HTTP сервер запущен: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Не удалось запустить HTTP сервер: %v\n", err)
	}
}
