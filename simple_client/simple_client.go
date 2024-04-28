package main

import (
	"bufio"
	"fmt"
	"net/http"
)

func main() {
	resp, err := http.Get("http://localhost:50051/")
	if err != nil {
		panic(err)
	} else if resp != nil {
		fmt.Println("Status: ", resp.Status)
		defer resp.Body.Close()
	}

	data := bufio.NewScanner(resp.Body)

	for data.Scan() {
		fmt.Println(data.Text())
	}

	if err := data.Err(); err != nil {
		panic(err)
	}
}
