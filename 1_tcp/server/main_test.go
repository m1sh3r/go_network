package main

import (
	"bufio"
	"net"
	"strings"
	"sync/atomic"
	"testing"
)

func TestHandleConnection(t *testing.T) {
	atomic.StoreUint64(&messageCount, 0)

	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	go handleConnection(serverConn)

	if _, err := clientConn.Write([]byte("hello\n")); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	response, err := bufio.NewReader(clientConn).ReadString('\n')
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if response != "Эхо hello\n" {
		t.Errorf("got %q, want %q", response, "Эхо hello\\n")
	}

	count := atomic.LoadUint64(&messageCount)
	if count != 1 {
		t.Errorf("got %d, want %d", count, 1)
	}
}

func TestHandleStats(t *testing.T) {
	atomic.StoreUint64(&messageCount, 7)

	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	go handleStats(serverConn)

	response, err := bufio.NewReader(clientConn).ReadString('\n')
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if !strings.Contains(response, "7") {
		t.Errorf("stats response does not contain count: %q", response)
	}
}
