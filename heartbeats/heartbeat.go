package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func SendHeartbeat(mainAddr string) {
	for {
		// Construct the URL for the heartbeat endpoint
		endpoint := fmt.Sprintf("http://%s/heartbeat", mainAddr)

		// Construct the POST data: empty data
		postData := url.Values{}

		// Send heartbeat POST request to master
		_, err := http.PostForm(endpoint, postData)
		if err != nil {
			log.Println("Error sending heartbeat:", err)
		} else {
			log.Println("Sent heartbeat to main")
		}

		time.Sleep(5 * time.Second) // Send heartbeat every 5 seconds
	}
}

func main() {
	// main address
	// If using Google Cloud, change this variable to be the address of the main server
	mainAddr := "localhost:8080"

	// Start sending heartbeats
	SendHeartbeat(mainAddr)
}
