package main

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"
	"web_main/consistent_hash"
)

type Main struct {
	mainPort      int
	nodeAddresses []string
	nodeMap       map[string]consistent_hash.ServerNode
}

type HotKeyEntry struct {
	Average         float64
	PastTimeRequest int64
}

type HotKeys struct {
	KeyMap map[string]HotKeyEntry
	mutex  sync.RWMutex
}

func Keys() *HotKeys {
	return &HotKeys{
		KeyMap: make(map[string]HotKeyEntry),
	}
}

// Get retrieves a cached entry by key
func (c *HotKeys) Get(url string) (HotKeyEntry, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, ok := c.KeyMap[url]
	return entry, ok
}

// Set inserts or updates a cached entry
func (hk *HotKeys) Set(url string, entry HotKeyEntry) {
	hk.mutex.Lock()
	defer hk.mutex.Unlock()
	hk.KeyMap[url] = entry
}

func NewMain(mainPort int, nodeList []consistent_hash.ServerNode) *Main {
	main := Main{mainPort: mainPort}

	nodeMap := make(map[string]consistent_hash.ServerNode)
	for _, node := range nodeList {
		nodeMap[node.IP] = node
		main.nodeAddresses = append(main.nodeAddresses, node.IP)
	}
	main.nodeMap = nodeMap

	return &main
}

func (main Main) updateNodeTimestamps(node string, w http.ResponseWriter) {
	nodeData, exists := main.nodeMap[node]
	if !exists {
		http.Error(w, "Node does not exist", http.StatusBadRequest)
		return
	}
	nodeData.Timestamp = time.Now()
	main.nodeMap[node] = nodeData
}

func (main Main) processPostRequest(w http.ResponseWriter, r *http.Request) {
	// TODO: change this later based upon IP address
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get the node address from the form data
	node := r.Form.Get("node")
	if node == "" {
		http.Error(w, "Node parameter is missing", http.StatusBadRequest)
		return
	}

	// Process the heartbeat (for example, you can log it)
	fmt.Printf("Received heartbeat from node %s\n", node)
	main.updateNodeTimestamps(node, w)

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
}

func (main Main) serve() {
	consistentHash := consistent_hash.NewTrie(main.nodeMap)

	// TODO figure out best threshold / k value
	threshhold := 3.0
	k := 0.5

	hotUrls := Keys()

	// Start the heartbeat server
	http.Handle("/heartbeat", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		main.processPostRequest(w, r)
	}))

	// Start the main server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Unix()

		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
			return
		}

		ip := ""
		// Find the IP address of the node that will serve the URL if url is hot
		value, exists := hotUrls.Get(url)
		if exists {
			if value.Average >= threshhold {
				fmt.Println("Threshold reached, randomly dispersing.")
				ip = main.nodeAddresses[rand.Intn(len(main.nodeAddresses))]
			} else {
				ip = consistentHash.Search(url)
			}

			if value.PastTimeRequest == now {
				fmt.Println("Same request in current second, calculating moving average.")
				hotUrls.Set(url, HotKeyEntry{
					Average:         value.Average + 1,
					PastTimeRequest: value.PastTimeRequest,
				})
			} else {
				fmt.Println("Moving average update for time of different second.")

				seconds := (float64)(now - value.PastTimeRequest)
				newAverage := value.Average*math.Pow(k, seconds) + 1
				hotUrls.Set(url, HotKeyEntry{
					Average:         newAverage,
					PastTimeRequest: now,
				})
			}
		} else {
			fmt.Println("Starting entry of moving average.")
			ip = consistentHash.Search(url)
			hotUrls.Set(url, HotKeyEntry{
				Average:         1,
				PastTimeRequest: now,
			})
		}

		for time.Since(main.nodeMap[ip].Timestamp) > 15*time.Second {
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
	// Initialize for time.Now() + 60 seconds to allow for starting everything up
	timestamp := time.Now().Add(60 * time.Second)
	nodeList := []consistent_hash.ServerNode{{IP: "localhost", Timestamp: timestamp, Replicas: 3}, {IP: "10.30.147.20", Timestamp: timestamp, Replicas: 3}}
	main := NewMain(8080, nodeList)
	main.serve()
}
