// package main

// import (
// 	"fmt"
// 	"hash/fnv"
// )

// func hashString(s string) uint32 {
// 	h := fnv.New32a()
// 	h.Write([]byte(s))
// 	return h.Sum32()
// }

// func main() {
// 	urls := []string{"www.google.com", "www.facebook.com", "www.youtube.com", "www.amazon.com", "www.netflix.com"}
// 	for _, url := range urls {
// 		fmt.Printf("URL: %v, Hash: %v\n", url, hashString(url)%256)
// 	}
// }
