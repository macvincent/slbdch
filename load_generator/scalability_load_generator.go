package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

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
func GenerateRequests(filePath string) []string {

	urlFrequencies, err := ReadCSV(filePath)
	if err != nil {
		log.Fatalf("Error reading CSV file: %v\n", err)
	}

	var requests []string
	for url, freq := range urlFrequencies {
		for i := 0; i < freq; i++ {
			requests = append(requests, url)
		}
	}
	// Shuffle the requests
	rand.Seed(1.0)
	rand.Shuffle(len(requests), func(i, j int) { requests[i], requests[j] = requests[j], requests[i] })

	return requests
}

func RequesWorker(requests []string, masterNode string, wg *sync.WaitGroup, outputThroughput int) {
	wokerWaitGroup := sync.WaitGroup{}
	wokerWaitGroup.Add(len(requests))
	ticker := time.NewTicker(time.Second / time.Duration(outputThroughput))
	failedRequests := 0
	defer ticker.Stop()

	for _, url := range requests {
		<-ticker.C
		go func(url string) {
			defer wokerWaitGroup.Done()
			defer wg.Done()
			url = fmt.Sprintf("%s/?url=%s", masterNode, url)
			resp, err := http.Get(url)
			if err != nil {
				failedRequests++
				log.Printf("Failed request: %v\n", err)
				return
			}
			resp.Body.Close()
			time.Sleep(1 * time.Millisecond)
		}(url)
	}
	wokerWaitGroup.Wait()
	log.Printf("Worker done. Failed requests: %d\n", failedRequests)
}

func main() {
	// Define the master node address
	masterNode := "http://localhost:8080"

	filePath := "load_generator/urlFrequencies.csv"
	requests := GenerateRequests(filePath)

	outputThroughput := 100 // requests per second
	numWorkers := 10
	if outputThroughput%numWorkers != 0 || outputThroughput*numWorkers < 0 {
		log.Fatalf("Invalid outputThroughput or numWorkers\n")
	}

	var wg sync.WaitGroup
	wg.Add(len(requests))

	ticker := time.NewTicker(time.Second / time.Duration(outputThroughput))
	defer ticker.Stop()

	startTime := time.Now()
	log.Printf("Starting load test\n")
	for i := 0; i < numWorkers; i++ {
		startIndex := i * len(requests) / numWorkers
		endIndex := (i + 1) * len(requests) / numWorkers
		go RequesWorker(requests[startIndex:endIndex], masterNode, &wg, outputThroughput/numWorkers)
	}
	log.Printf("All requests sent\n")
	wg.Wait()
	elapsedTime := time.Since(startTime)

	log.Printf("Total requests: %d\n", len(requests))
	log.Printf("Elapsed time: %s\n", elapsedTime)
	log.Printf("Request Throughput: %d requests per second\n", outputThroughput)
	log.Printf("Response Throughput: %.2f requests per second\n", float64(len(requests))/elapsedTime.Seconds())
}
