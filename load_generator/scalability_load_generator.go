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

func main() {
	// Define the master node address
	masterNode := "http://localhost:8080"

	filePath := "load_generator/urlFrequencies.csv"

	requests := GenerateRequests(filePath)

	var wg sync.WaitGroup
	wg.Add(len(requests))

	outputThroughput := 100.0 // requests per second
	ticker := time.NewTicker(time.Second / time.Duration(outputThroughput))
	defer ticker.Stop()

	startTime := time.Now()
	successCounter := 0
	failCounter := 0

	// TODO: One worker might prove to be a bottleneck. Consider using multiple workers.
	log.Printf("Starting load test\n")
	for _, url := range requests {
		<-ticker.C
		go func(url string) {
			defer wg.Done()
			url = fmt.Sprintf("%s/?url=%s", masterNode, url)
			resp, err := http.Get(url)
			if err != nil {
				failCounter++
				log.Printf("Failed request: %v\n", err)
				return
			}
			resp.Body.Close()
			successCounter++
		}(url)
	}
	log.Printf("All requests sent\n")
	wg.Wait()
	elapsedTime := time.Since(startTime)

	log.Printf("Total requests: %d\n", len(requests))
	log.Printf("Successful requests: %d\n", successCounter)
	log.Printf("Failed requests: %d\n", failCounter)
	log.Printf("Elapsed time: %s\n", elapsedTime)
	log.Printf("Request Throughput: %.2f requests per second\n", outputThroughput)
	log.Printf("Response Throughput: %.2f requests per second\n", float64(successCounter)/elapsedTime.Seconds())
}
