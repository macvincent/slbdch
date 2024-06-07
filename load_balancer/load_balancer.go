package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	mux sync.RWMutex
)

type MasterNode struct {
	IP        string
	Timestamp time.Time
}

type Main struct {
	mainPort      int
	nodeAddresses []string
	nodeMap       map[string]MasterNode
}

func NewMain(mainPort int, nodeList []MasterNode) *Main {
	main := Main{mainPort: mainPort}

	nodeMap := make(map[string]MasterNode)
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

func (main *Main) processInsert(w http.ResponseWriter, r *http.Request) {
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
	if new_node_ip_address == "" {
		http.Error(w, "Some parameter is missing", http.StatusBadRequest)
		return
	}

	mux.Lock()
	main.nodeMap[new_node_ip_address] = MasterNode{IP: new_node_ip_address, Timestamp: time.Now()}
	main.nodeAddresses = append(main.nodeAddresses, new_node_ip_address)
	mux.Unlock()

	fmt.Printf("Inserted new node %s\n", ip_address)
	// Respond with a success message
	w.WriteHeader(http.StatusOK)
}

func (main *Main) processDelete(w http.ResponseWriter, r *http.Request) {
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

	main.DeleteNode(remove_ip_address)

	// Process the heartbeat (for example, you can log it)
	fmt.Printf("Deleted node %s\n", ip_address)

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
}

func (main *Main) DeleteNode(ip string) {
	mux.Lock()
	delete(main.nodeMap, ip)
	for i, node := range main.nodeAddresses {
		if node == ip {
			main.nodeAddresses = append(main.nodeAddresses[:i], main.nodeAddresses[i+1:]...)
			break
		}
	}
	mux.Unlock()
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

		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
			return
		}

		// Find the IP address of the server to send the request to randomly from node address list
		if main.nodeAddresses == nil || len(main.nodeAddresses) == 0 {
			http.Error(w, "No nodes available", http.StatusInternalServerError)
			return
		}

		ip := main.nodeAddresses[rand.Intn(len(main.nodeAddresses))]
		for time.Since(main.nodeMap[ip].Timestamp) > 15*time.Second {
			main.DeleteNode(ip)
			if len(main.nodeAddresses) == 0 {
				http.Error(w, "No nodes available", http.StatusInternalServerError)
				return
			}
			ip = main.nodeAddresses[rand.Intn(len(main.nodeAddresses))]
		}

		// Send request to found ip address
		http.Redirect(w, r, fmt.Sprintf("http://%v:8080?url=%v", ip, url), http.StatusTemporaryRedirect)
	})

	serveAddr := fmt.Sprintf(":%d", main.mainPort)
	fmt.Println("Server started on ", serveAddr)
	http.ListenAndServe(serveAddr, nil)
}

func main() {
	// Initialize for time.Now() + 60 seconds to allow for starting everything up
	timestamp := time.Now().Add(60 * time.Second)
	nodeList := []MasterNode{{IP: "localhost", Timestamp: timestamp}}
	main := NewMain(5037, nodeList)
	main.serve()
}
