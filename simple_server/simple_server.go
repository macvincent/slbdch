package main

import (
	"fmt"
	"log"
	"net/http"
)

func defaultResponse(writer http.ResponseWriter, request *http.Request) {
	log.Println("Default request received...")
	fmt.Fprintf(writer, "Hello there!")
}

func requestHeaders(writer http.ResponseWriter, request *http.Request) {
	log.Println("Header request received...")
	fmt.Fprintf(writer, "Name: Header\n")
	for name, headers := range request.Header {
		for _, h := range headers {
			fmt.Fprintf(writer, "%v: %v\n", name, h)
		}
	}
}

func main() {
	http.HandleFunc("/", defaultResponse)
	http.HandleFunc("/headers", requestHeaders)

	fmt.Println("Server is listening on port 50051...")
	http.ListenAndServe(":50051", nil)
}
