package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	latencyFile *os.File
	fileMutex   sync.Mutex
)

func recordLatency(latency time.Duration) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	_, err := latencyFile.WriteString(fmt.Sprintf("%v\n", latency.Nanoseconds()))
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func init() {
	var err error
	latencyFile, err = os.OpenFile("load_tester_latencies.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
}

// ReadCSV reads a CSV file and returns a map of URLs and their frequencies
func ReadCSV(filePath string) (map[string]int, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	urlFrequencies := make(map[string]int)
	for _, record := range records[1:] {
		if len(record) != 2 {
			return nil, fmt.Errorf("invalid CSV format")
		}
		url := record[0]
		freq, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}
		urlFrequencies[url] = freq
	}
	return urlFrequencies, nil
}

// GenerateRequests generates a slice of URLs based on their frequencies
func GenerateRequests(urlFrequencies map[string]int) []string {
	var requests []string
	for url, freq := range urlFrequencies {
		for i := 0; i < freq; i++ {
			requests = append(requests, url)
		}
	}
	return requests
}

func worker(client *http.Client, masterNode string, jobs <-chan string, wg *sync.WaitGroup, successCounter *int, failCounter *int, mutex *sync.Mutex, w int) {
	defer wg.Done()

	for rawURL := range jobs {

		start := time.Now()
		// Construct the URL to be sent to the master node
		masterURL := fmt.Sprintf("%s/?url=%s", masterNode, url.QueryEscape(rawURL))
		resp, err := client.Get(masterURL)

		if err != nil {
			mutex.Lock()
			*failCounter++
			mutex.Unlock()
			log.Printf("Failed to send request to %s: %v\n", masterURL, err)
			continue
		}

		resp.Body.Close()
		mutex.Lock()
		*successCounter++
		mutex.Unlock()
		end := time.Now()
		latency := end.Sub(start)
		recordLatency(latency)
		log.Printf("Successfully sent request to %s by worker %d\n", masterURL, w)
	}
}

func main() {
	// Define the master node address
	masterNode := "http://localhost:8080"

	// Path to the CSV file
	filePath := "load_generator/urlFrequencies.csv"

	// Number of parallel workers
	numWorkers := 1

	urlFrequencies, err := ReadCSV(filePath)
	if err != nil {
		log.Fatalf("Error reading CSV file: %v\n", err)
	}

	requests := GenerateRequests(urlFrequencies)
	rand.Seed(1)
	rand.Shuffle(len(requests), func(i, j int) { requests[i], requests[j] = requests[j], requests[i] })

	jobs := make(chan string, len(requests))
	var wg sync.WaitGroup
	var successCounter, failCounter int
	var mutex sync.Mutex
	client := &http.Client{}

	startTime := time.Now()

	// Start workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go worker(client, masterNode, jobs, &wg, &successCounter, &failCounter, &mutex, w)
	}

	// Send jobs to the workers
	for _, rawURL := range requests {
		jobs <- rawURL
	}
	close(jobs)

	wg.Wait()
	elapsedTime := time.Since(startTime)

	log.Printf("Total requests: %d\n", len(requests))
	log.Printf("Successful requests: %d\n", successCounter)
	log.Printf("Failed requests: %d\n", failCounter)
	log.Printf("Elapsed time: %s\n", elapsedTime)
}
