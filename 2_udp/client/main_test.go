package main

import (
	"net"
	"testing"
	"time"
)

func TestExchangeMessage(t *testing.T) {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	defer serverConn.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		buffer := make([]byte, 1024)
		n, clientAddr, err := serverConn.ReadFromUDP(buffer)
		if err != nil {
			return
		}
		if n == 0 {
			return
		}
		_, _ = serverConn.WriteToUDP([]byte("Привет!"), clientAddr)
	}()

	clientConn, err := net.DialUDP("udp", nil, serverConn.LocalAddr().(*net.UDPAddr))
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer clientConn.Close()

	response, err := exchangeMessage(clientConn, "hello")
	if err != nil {
		t.Fatalf("exchange failed: %v", err)
	}

	if response != "Привет!" {
		t.Errorf("got %q, want %q", response, "Привет!")
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("server goroutine timeout")
	}
}
