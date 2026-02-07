package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8081")
	if err != nil {
		fmt.Printf("Не удалось определить адрес UDP: %v\n", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Printf("Не удалось подключиться по UDP: %v\n", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Подключено к серверу. Введите сообщение и нажмите Enter. Для выхода наберите exit.")
	for {
		fmt.Print("Сообщение: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Не удалось прочитать ввод: %v\n", err)
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.EqualFold(line, "exit") {
			return
		}

		if _, err := conn.Write([]byte(line)); err != nil {
			fmt.Printf("Не удалось отправить сообщение: %v\n", err)
			return
		}

		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			fmt.Printf("Не удалось установить таймаут: %v\n", err)
			return
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Превышено время ожидания ответа")
				continue
			}
			fmt.Printf("Не удалось прочитать ответ: %v\n", err)
			return
		}

		fmt.Printf("%s\n", string(buf[:n]))
	}
}
