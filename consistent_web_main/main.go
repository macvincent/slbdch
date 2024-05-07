package main

import (
	"fmt"
	"io"
	"net/http"
	"web_main/consistent_hash"
)

func main() {
	// TODO: Update this with local web cache port numbers
	nodeAddresses := []string{"56085", "56086", "56087", "56088", "56089"}
	replicaPerNode := 10
	consistentHash := consistent_hash.NewTrie(nodeAddresses, replicaPerNode)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	fmt.Println("Server started on :5050")
	http.ListenAndServe(":5050", nil)
}
