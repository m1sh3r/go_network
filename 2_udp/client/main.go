package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func exchangeMessage(conn *net.UDPConn, line string) (string, error) {
	if strings.TrimSpace(line) == "" {
		return "", errors.New("пустое сообщение")
	}

	if _, err := conn.Write([]byte(line)); err != nil {
		return "", fmt.Errorf("не удалось отправить сообщение: %w", err)
	}

	if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return "", fmt.Errorf("не удалось установить таймаут: %w", err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать ответ: %w", err)
	}

	return string(buf[:n]), nil
}

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

		response, err := exchangeMessage(conn, line)
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				fmt.Println("Превышено время ожидания ответа")
				continue
			}
			fmt.Printf("%v\n", err)
			return
		}

		fmt.Printf("%s\n", response)
	}
}
