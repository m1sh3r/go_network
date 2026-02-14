package main

import (
	"fmt"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:8081")
	if err != nil {
		fmt.Printf("Не удалось определить адрес UDP: %v\n", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("Не удалось запустить UDP сервер: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("UDP сервер слушает: localhost:8081")

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Не удалось прочитать UDP: %v\n", err)
			continue
		}
		msg := string(buffer[:n])
		fmt.Printf("Сообщение от %s: %s", clientAddr, msg)

		// Отправляем ответ
		if _, err := conn.WriteToUDP([]byte("Привет!"), clientAddr); err != nil {
			fmt.Printf("Не удалось отправить ответ: %v\n", err)
		}
	}
}
