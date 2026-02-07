package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Не удалось подключиться: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Подключено к серверу. Введите сообщение и нажмите Enter.")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Сообщение: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Printf("Не удалось прочитать ввод: %v\n", err)
			}
			return
		}
		text := scanner.Text()
		if _, err := conn.Write([]byte(text + "\n")); err != nil {
			fmt.Printf("Не удалось отправить сообщение: %v\n", err)
			return
		}

		response := make([]byte, 1024)
		n, err := conn.Read(response)
		if err != nil {
			fmt.Printf("Не удалось прочитать ответ: %v\n", err)
			return
		}
		fmt.Print(string(response[:n]))
	}
}
