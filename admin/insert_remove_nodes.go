package admin

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func SendInsertNodeCommand(mainAddr string, ip_address string, replica_count string) {
	// Construct the URL for the heartbeat endpoint
	endpoint := fmt.Sprintf("http://%s/insert", mainAddr)

	// Construct the POST data: empty data
	postData := url.Values{}
	postData.Set("ip_address", ip_address)
	postData.Set("replica_count", replica_count)

	// Send heartbeat POST request to master
	_, err := http.PostForm(endpoint, postData)
	if err != nil {
		log.Println("Error sending insert:", err)
	} else {
		log.Println("Inserted new node")
	}
}

func SendRemoveNodeCommand(mainAddr string, ip_address string) {
	// Just in case we want to remove a given node immediately

	// Construct the URL for the heartbeat endpoint
	endpoint := fmt.Sprintf("http://%s/delete", mainAddr)

	// Construct the POST data: empty data
	postData := url.Values{}
	postData.Set("ip_address", ip_address)

	// Send heartbeat POST request to master
	_, err := http.PostForm(endpoint, postData)
	if err != nil {
		log.Println("Error sending delete:", err)
	} else {
		log.Println("Deleted node")
	}
}

func main() {
	// Master address
	masterAddr := "localhost:5050"

	if os.Args[1] == "insert" {
		SendInsertNodeCommand(masterAddr, os.Args[2], os.Args[3])
	} else if os.Args[1] == "remove" {
		SendRemoveNodeCommand(masterAddr, os.Args[2])
	}
}
