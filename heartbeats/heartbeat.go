package main

import (
	"fmt"
	"net"
	"time"
)

func sendHeartbeat(conn net.Conn) {
	for {
		_, err := conn.Write([]byte("heartbeat"))
		if err != nil {
			fmt.Println("Error sending heartbeat:", err)
			return
		}
		time.Sleep(5 * time.Second) // Send heartbeat every 5 seconds
	}
}

func main() {
	// Master address
	masterAddr := "localhost:8080"

	// Connect to master
	conn, err := net.Dial("tcp", masterAddr)
	if err != nil {
		fmt.Println("Error connecting to master:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to master at", masterAddr)

	// Start sending heartbeats
	sendHeartbeat(conn)
}
