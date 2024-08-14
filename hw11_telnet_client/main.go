package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go-telnet [--timeout=duration] host port")
		os.Exit(1)
	}

	timeout := 10 * time.Second
	host := os.Args[len(os.Args)-2]
	port := os.Args[len(os.Args)-1]

	for i := 1; i < len(os.Args)-2; i++ {
		if os.Args[i] == "--timeout" {
			if i+1 < len(os.Args)-2 {
				var err error
				timeout, err = time.ParseDuration(os.Args[i+1])
				if err != nil {
					fmt.Println("Invalid timeout format:", err)
					os.Exit(1)
				}
			}
		}
	}

	address := net.JoinHostPort(host, port)
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := client.Connect()
	if err != nil {
		fmt.Println("Connection error:", err)
		os.Exit(1)
	}
	defer client.Close()

	done := make(chan struct{})
	go func() {
		for {
			err := client.Send()
			if err != nil {
				fmt.Println("Send error:", err)
				close(done)
				return
			}
		}
	}()

	go func() {
		for {
			err := client.Receive()
			if err != nil {
				fmt.Println("Receive error:", err)
				close(done)
				return
			}
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
		fmt.Println("Terminated by signal")
	case <-done:
		fmt.Println("Connection closed")
	case <-ctx.Done():
		fmt.Println("Context timeout or canceled")
	}
}
