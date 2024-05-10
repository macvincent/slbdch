package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func sendHeartbeat(masterAddr string, port string) {
	for {
		// Construct the URL for the heartbeat endpoint
		endpoint := fmt.Sprintf("http://%s/heartbeat", masterAddr)

		// Construct the POST data
		postData := url.Values{}
		postData.Set("port", port)

		// Send heartbeat POST request to master
		_, err := http.PostForm(endpoint, postData)
		if err != nil {
			fmt.Println("Error sending heartbeat:", err)
		} else {
			fmt.Println("Sent heartbeat to master from port:", port)
		}

		time.Sleep(5 * time.Second) // Send heartbeat every 5 seconds
	}
}

func main() {
	// Master address
	masterAddr := "localhost:5050"

	// Port of this server. TODO change this later
	port := "12345"

	// Start sending heartbeats
	sendHeartbeat(masterAddr, port)
}
