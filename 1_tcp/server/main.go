package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync/atomic"
)

var messageCount uint64

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Клиент отключился")
			} else {
				fmt.Printf("Не удалось прочитать сообщение: %v\n", err)
			}
			return
		}
		message = strings.TrimSpace(message)
		atomic.AddUint64(&messageCount, 1)
		fmt.Printf("Получено сообщение: %s\n", message)
		if _, err := conn.Write([]byte("Эхо " + message + "\n")); err != nil {
			fmt.Printf("Не удалось отправить ответ: %v\n", err)
			return
		}
	}
}

func handleStats(conn net.Conn) {
	defer conn.Close()

	count := atomic.LoadUint64(&messageCount)
	if _, err := fmt.Fprintf(conn, "Сообщений: %d\n", count); err != nil {
		fmt.Printf("Не удалось отправить статистику: %v\n", err)
	}
}

func serveStats(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Не удалось запустить сервер статистики: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Сервер статистики слушает: %s\n", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Не удалось принять соединение статистики: %v\n", err)
			continue
		}
		go handleStats(conn)
	}
}

func main() {
	go serveStats("localhost:8082")

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Не удалось запустить TCP сервер: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP сервер слушает: localhost:8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Не удалось принять соединение: %v\n", err)
			continue
		}
		go handleConnection(conn) // каждый клиент — в отдельной горутине
	}
}
