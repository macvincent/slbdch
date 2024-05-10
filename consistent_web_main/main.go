package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
	"web_main/consistent_hash"
)

type Master struct {
	masterPort       int
	nodeAddresses    []string
	replicaPerNode   int
	portTimestampMap map[string]time.Time
}

func (master Master) processPostRequest(w http.ResponseWriter, r *http.Request) {
	// TODO: change this later based upon IP address
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get the port from the form data
	port := r.Form.Get("port")
	if port == "" {
		http.Error(w, "Port parameter is missing", http.StatusBadRequest)
		return
	}

	// Process the heartbeat (for example, you can log it)
	fmt.Printf("Received heartbeat from port %s\n", port)
	master.portTimestampMap[port] = time.Now()

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
}

func (master Master) serve() {
	consistentHash := consistent_hash.NewTrie(master.nodeAddresses, master.replicaPerNode)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			master.processPostRequest(w, r)
		}
		if r.Method == http.MethodGet {
			url := r.URL.Query().Get("url")
			if url == "" {
				http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
				return
			}

			// Find the IP address of the node that will serve the URL
			ip := consistentHash.Search(url)
			// Fetch content from the web
			resp, err := http.Get(fmt.Sprintf("http://localhost:%v//%v?url=%v", ip, ip, url))

			if err != nil {
				http.Error(w, fmt.Sprintf("Error fetching URL: %v", err), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			content, err := io.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error reading response body: %v", err), http.StatusInternalServerError)
				return
			}

			// Serve fetched content
			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			w.Write(content)
		}
	})

	serveAddr := fmt.Sprintf(":%d", master.masterPort)
	fmt.Println("Server started on ", serveAddr)
	http.ListenAndServe(serveAddr, nil)
}

func main() {
	// TODO: Update this with local web cache port numbers
	portTimestampMap := make(map[string]time.Time)
	nodeAddresses := []string{"58535", "58536", "58537", "58538", "58539"}

	// Initialize for time.Now() + 10 seconds to allow for starting
	// everything up
	timestamp := time.Now().Add(10 * time.Second)
	for _, port := range nodeAddresses {
		portTimestampMap[port] = timestamp
	}

	master := Master{
		masterPort:       5050,
		nodeAddresses:    nodeAddresses,
		replicaPerNode:   10,
		portTimestampMap: portTimestampMap,
	}
	master.serve()
}
