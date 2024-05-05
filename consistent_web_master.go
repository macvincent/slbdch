package main

import (
	"final_project/consistent_hash"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Create a new consistent hash with 5 nodes and 10 replicas per node
	nodeAddresses := []string{"8.8.8.8", "1.1.1.1", "208.67.222.222", "208.67.220.220", "9.9.9.9"}
	replicaPerNode := 10
	consistentHash := consistent_hash.NewConsistentHash(nodeAddresses, replicaPerNode)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
			return
		}

		// Find the IP address of the node that will serve the URL
		ip := consistentHash.ValueLookup(url)
		// Fetch content from the web
		resp, err := http.Get(fmt.Sprintf("http://%v?url=%v", ip, url))

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

	fmt.Println("Server started on :5050")
	http.ListenAndServe(":5050", nil)
}
