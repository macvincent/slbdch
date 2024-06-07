package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type RequestData struct {
	sentTime     time.Time
	responseTime time.Time
	host         string
}

var wg sync.WaitGroup

func sendRequest(client *http.Client, url string, results chan<- RequestData) {
	defer wg.Done()
	reqData := RequestData{sentTime: time.Now()}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	reqData.responseTime = time.Now()
	reqData.host = resp.Request.URL.Host
	results <- reqData
}

func main() {
	masterAddr := "http://34.125.70.184:8080"
	rawURL := "https://www.google.com"
	requestUrl := fmt.Sprintf("%s/?url=%s", masterAddr, url.QueryEscape(rawURL))
	requestsPerSecond := []int{} // Example input
	results := make(chan RequestData, 10000)

	client := &http.Client{}

	startTime := time.Now()

	go func() {
		for i, rps := range requestsPerSecond {
			secondStart := startTime.Add(time.Duration(i) * time.Second)
			ticker := time.NewTicker(time.Second / time.Duration(rps))
			defer ticker.Stop()

			for j := 0; j < rps; j++ {
				<-ticker.C
				wg.Add(1)
				go sendRequest(client, requestUrl, results)
			}

			// Wait until the start of the next second
			time.Sleep(time.Until(secondStart.Add(time.Second)))
		}
		wg.Wait()
		close(results)
	}()

	var sentPerSecond = make(map[int]int)
	var responsesPerSecond = make(map[string]map[int]int)

	for result := range results {
		sentSecond := int(result.sentTime.Sub(startTime).Seconds())
		responseSecond := int(result.responseTime.Sub(startTime).Seconds())
		sentPerSecond[sentSecond]++

		if responsesPerSecond[result.host] == nil {
			responsesPerSecond[result.host] = make(map[int]int)
		}
		responsesPerSecond[result.host][responseSecond]++
	}

	fmt.Println("Requests sent per second:", sentPerSecond)
	for host, perSecondMap := range responsesPerSecond {
		fmt.Printf("Responses received per second for host %s: %v\n", host, perSecondMap)
	}
}
