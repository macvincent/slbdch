package main

import (
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	"web_main/consistent_hash"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	latencyFile     *os.File
	fileMutex       sync.Mutex
	latencyFileName string
)

type Main struct {
	mainPort       int
	nodeAddresses  []string
	nodeMap        map[string]consistent_hash.ServerNode
	consistentHash *consistent_hash.Cycle
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

	main.consistentHash = consistent_hash.NewConsistentHash(nodeMap)

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

func (main Main) processHeartbeat(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	// Get the port from the form data
	ip_address, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil || ip_address == "" {
		http.Error(w, "Cannot get IP address", http.StatusBadRequest)
		return
	}
	if ip_address == "::1" {
		ip_address = "localhost" // localhost or 127.0.0.1 is equivalent to ::1
	}

	// Process the heartbeat (for example, you can log it)
	fmt.Printf("Received heartbeat from node %s\n", ip_address)
	main.updateNodeTimestamps(ip_address, w)

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
}

func (main Main) processInsert(w http.ResponseWriter, r *http.Request) {
	// Get the port from the form data
	ip_address, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil || ip_address != "::1" {
		http.Error(w, "Cannot identify valid host", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	new_node_ip_address := r.Form.Get("ip_address")
	new_node_replica_count := r.Form.Get("replica_count")
	if new_node_ip_address == "" || new_node_replica_count == "" {
		http.Error(w, "Some parameter is missing", http.StatusBadRequest)
		return
	}
	new_node_replica_count_int, err := strconv.Atoi(new_node_replica_count)
	if err != nil {
		http.Error(w, "Error parsing replica count", http.StatusBadRequest)
		return
	}

	main.consistentHash.InsertNode(new_node_ip_address, new_node_replica_count_int)

	// Process the heartbeat (for example, you can log it)
	fmt.Printf("Inserted new node %s\n", ip_address)

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
}

func (main Main) processDelete(w http.ResponseWriter, r *http.Request) {
	// Get the port from the form data
	ip_address, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil || ip_address != "::1" {
		http.Error(w, "Cannot identify valid host", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	remove_ip_address := r.Form.Get("ip_address")
	if remove_ip_address == "" {
		http.Error(w, "Some parameter is missing", http.StatusBadRequest)
		return
	}

	main.consistentHash.DeleteNode(remove_ip_address)

	// Process the heartbeat (for example, you can log it)
	fmt.Printf("Deleted node %s\n", ip_address)

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
}

func recordLatency(latency time.Duration) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	_, err := latencyFile.WriteString(fmt.Sprintf("%v\n", latency.Nanoseconds()))
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func saveAndCloseFile() {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	err := latencyFile.Close()
	if err != nil {
		fmt.Println("Error closing file:", err)
		return
	}
	// Reopen the file for further writing
	latencyFile, err = os.OpenFile(latencyFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error reopening file:", err)
		os.Exit(1)
	}
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

	// TODO figure out best threshold / k value
	threshhold := 10000.0
	k := 0.5

	hotUrls := Keys()

	// Start the heartbeat server
	http.Handle("/heartbeat", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		main.processHeartbeat(w, r)
	}))

	http.Handle("/insert", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		main.processInsert(w, r)
	}))

	http.Handle("/delete", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		main.processDelete(w, r)
	}))

	// Start the main server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Unix()

		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
			return
		}
		if url == "done" {
			saveAndCloseFile()
		}

		ip := ""
		// Find the IP address of the node that will serve the URL if url is hot
		value, exists := hotUrls.Get(url)
		if exists {
			if value.Average >= threshhold {
				//logger.Info("Threshold reached, randomly dispersing.")
				ip = main.nodeAddresses[rand.Intn(len(main.nodeAddresses))]
			} else {
				start_time := time.Now()
				ip = main.consistentHash.ValueLookup(url)
				end_time := time.Now()
				latency := end_time.Sub(start_time)
				recordLatency(latency)
			}

			if value.PastTimeRequest == now {
				logger.Info("Same request in current second, calculating moving average.")
				hotUrls.Set(url, HotKeyEntry{
					Average:         value.Average + 1,
					PastTimeRequest: value.PastTimeRequest,
				})
			} else {
				//logger.Info("Moving average update for time of different second.")

				seconds := (float64)(now - value.PastTimeRequest)
				newAverage := value.Average*math.Pow(k, seconds) + 1
				hotUrls.Set(url, HotKeyEntry{
					Average:         newAverage,
					PastTimeRequest: now,
				})
			}
		} else {
			//logger.Info("Starting entry of moving average.")
			ip = main.consistentHash.ValueLookup(url)
			hotUrls.Set(url, HotKeyEntry{
				Average:         1,
				PastTimeRequest: now,
			})
		}

		for time.Since(main.nodeMap[ip].Timestamp) > 15*time.Second {
			main.consistentHash.DeleteNode(ip)
			ip = main.consistentHash.ValueLookup(url)
		}

		// Send request to found ip address
		http.Redirect(w, r, fmt.Sprintf("http://%v:5050?url=%v", ip, url), http.StatusTemporaryRedirect)
	})

	serveAddr := fmt.Sprintf(":%d", main.mainPort)
	fmt.Println("Server started on ", serveAddr)
	http.ListenAndServe(serveAddr, nil)
}

func init() {
	latencyFileName = fmt.Sprintf("latency_%v_replicas.csv", num_total_replicas)
	var err error
	latencyFile, err = os.OpenFile(latencyFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
}

var num_total_replicas int = 1

func main() {
	defer latencyFile.Close()
	runTests := false
	if runTests {
		consistent_hash.CycleMain()
		consistent_hash.KademliaMain()
	} else {
		// Initialize for time.Now() + 60 seconds to allow for starting everything up
		timestamp := time.Now().Add(60 * time.Second)
		// If using Google Cloud, change this variable to include the IPs of the servers you have created
		nodeList := []consistent_hash.ServerNode{
			{IP: "34.16.223.151", Timestamp: timestamp, Replicas: num_total_replicas}, // go-vm1
		}
		main := NewMain(8080, nodeList)
		main.serve()
	}
}
