package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func SendHeartbeat(mainAddr string, currAddress string) {
	for {
		// Construct the URL for the heartbeat endpoint
		endpoint := fmt.Sprintf("http://%s/heartbeat", mainAddr)

		// Construct the POST data
		postData := url.Values{}
		postData.Set("node", currAddress)

		// Send heartbeat POST request to master
		_, err := http.PostForm(endpoint, postData)
		if err != nil {
			log.Println("Error sending heartbeat:", err)
		} else {
			log.Println("Sent heartbeat to main from current address:", currAddress)
		}

		time.Sleep(5 * time.Second) // Send heartbeat every 5 seconds
	}
}

func main() {
	// Master address
	masterAddr := "localhost:5050"

	// Port of this server. TODO change this later
	currAddress := "58535"

	// Start sending heartbeats
	SendHeartbeat(masterAddr, currAddress)
}
