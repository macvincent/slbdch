package main

import (
	"fmt"
	"go.uber.org/zap"
    "go.uber.org/zap/zapcore"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"
	"web_main/consistent_hash"
)

type Main struct {
	mainPort         int
	nodeAddresses    []string
	replicaPerNode   int
	nodeTimestampMap map[string]time.Time
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
	// Create logger configuration with asynchronous logging enabled
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := cfg.Build()
	defer logger.Sync()


	consistentHash := consistent_hash.NewTrie(main.nodeAddresses, main.replicaPerNode)

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
				logger.Info("Threshold reached, randomly dispersing.")
				ip = main.nodeAddresses[rand.Intn(len(main.nodeAddresses))]
			} else {
				ip = consistentHash.Search(url)
			}

			if value.PastTimeRequest == now {
				logger.Info("Same request in current second, calculating moving average.")
				hotUrls.Set(url, HotKeyEntry{
					Average:         value.Average + 1,
					PastTimeRequest: value.PastTimeRequest,
				})
			} else {
				logger.Info("Moving average update for time of different second.")

				seconds := (float64)(now - value.PastTimeRequest)
				newAverage := value.Average*math.Pow(k, seconds) + 1
				hotUrls.Set(url, HotKeyEntry{
					Average:         newAverage,
					PastTimeRequest: now,
				})
			}
		} else {
			logger.Info("Starting entry of moving average.")
			ip = consistentHash.Search(url)
			hotUrls.Set(url, HotKeyEntry{
				Average:         1,
				PastTimeRequest: now,
			})
		}

		for time.Now().Sub(main.nodeTimestampMap[ip]) > 15*time.Second {
			consistentHash.DeleteNode(ip)
			ip = consistentHash.Search(url)
		}

		// Send request to found ip address
		http.Redirect(w, r, fmt.Sprintf("http://%v:8080?url=%v", ip, url), http.StatusTemporaryRedirect)
	})

	serveAddr := fmt.Sprintf(":%d", main.mainPort)
	fmt.Println("Server started on ", serveAddr)
	http.ListenAndServe(serveAddr, nil)
}

func main() {
	nodeTimestampMap := make(map[string]time.Time)
	nodeAddresses := []string{"localhost"}

	// Initialize for time.Now() + 500 seconds to allow for starting
	// everything up
	timestamp := time.Now().Add(500 * time.Second)
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
