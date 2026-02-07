package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("https://httpbin.org/json")
	if err != nil {
		fmt.Printf("Не удалось выполнить запрос: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Не удалось прочитать ответ: %v\n", err)
		return
	}

	fmt.Printf("Статус: %s\n", resp.Status)
	fmt.Printf("Тело ответа:\n%s\n", string(body))
}
