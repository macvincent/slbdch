package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
	"web_main/consistent_hash"
)

type Main struct {
	mainPort         int
	nodeAddresses    []string
	replicaPerNode   int
	nodeTimestampMap map[string]time.Time
}

func (main Main) processPostRequest(w http.ResponseWriter, r *http.Request) {
	// TODO: change this later based upon IP address
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get the port from the form data
	node := r.Form.Get("node")
	if node == "" {
		http.Error(w, "Node parameter is missing", http.StatusBadRequest)
		return
	}

	// Process the heartbeat (for example, you can log it)
	fmt.Printf("Received heartbeat from node %s\n", node)
	main.nodeTimestampMap[node] = time.Now()

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
}

func (main Main) serve() {
	consistentHash := consistent_hash.NewTrie(main.nodeAddresses, main.replicaPerNode)

	// Start the heartbeat server
	http.Handle("/heartbeat", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		main.processPostRequest(w, r)
	}))

	// Start the main server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
			return
		}

		// Find the IP address of the node that will serve the URL
		ip := consistentHash.Search(url)

		for time.Now().Sub(main.nodeTimestampMap[ip]) > 15*time.Second {
			consistentHash.DeleteNode(ip)
			ip = consistentHash.Search(url)
		}

		// Fetch content from the web
		resp, err := http.Get(fmt.Sprintf("http://%v:8080?url=%v", ip, url))

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
	})

	serveAddr := fmt.Sprintf(":%d", main.mainPort)
	fmt.Println("Server started on ", serveAddr)
	http.ListenAndServe(serveAddr, nil)
}

func main() {
	nodeTimestampMap := make(map[string]time.Time)
	nodeAddresses := []string{"34.125.246.49", "34.125.32.120"}

	// Initialize for time.Now() + 60 seconds to allow for starting
	// everything up
	timestamp := time.Now().Add(60 * time.Second)
	for _, nodeAddress := range nodeAddresses {
		nodeTimestampMap[nodeAddress] = timestamp
	}

	main := Main{
		mainPort:         5050,
		nodeAddresses:    nodeAddresses,
		replicaPerNode:   10,
		nodeTimestampMap: nodeTimestampMap,
	}
	main.serve()
}
