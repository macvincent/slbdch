package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// CacheMetrics represents metrics for the cache
type CacheMetrics struct {
	Hits         int
	RequestCount int
}

// CacheEntry represents an entry in the cache
type CacheEntry struct {
	Content     []byte
	ContentType string
	Expiration  time.Time
}

// Cache represents an in-memory cache
type Cache struct {
	entries map[string]CacheEntry
	mutex   sync.RWMutex
	metrics CacheMetrics
}

// NewCache creates a new instance of Cache
func NewCache() *Cache {
	return &Cache{
		entries: make(map[string]CacheEntry),
		metrics: CacheMetrics{0, 0},
	}
}

// Get retrieves a cached entry by key
func (c *Cache) Get(key string) (CacheEntry, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, ok := c.entries[key]

	if ok && time.Now().Before(entry.Expiration) {
		c.metrics.RequestCount++
		c.metrics.Hits++
	} else if ok {
		// Entry has expired
		return entry, false
	}

	return entry, ok
}

// Set inserts or updates a cached entry
func (c *Cache) Set(key string, entry CacheEntry) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.entries[key] = entry
	c.metrics.RequestCount++
}

func main() {
	cache := NewCache()

	httpAddr := flag.String("http", ":8080", "HTTP service address")

	log.Printf("HTTP service listening on %s", *httpAddr)

	// Expose metrics as JSON on /metrics
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		cache.mutex.RLock()
		defer cache.mutex.RUnlock()

		json := []byte(fmt.Sprintf(`{"hits": %d, "requests": %d}`, cache.metrics.Hits, cache.metrics.RequestCount))
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	})

	// Serve cached or fetched content
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
			return
		}

		if entry, ok := cache.Get(url); ok {
			// Serve cached content
			w.Header().Set("Content-Type", entry.ContentType)
			w.Write(entry.Content)
			log.Println("Served from cache:", url)
			return
		}

		// Fetch content from the web
		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching URL: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading response body: %v", err), http.StatusInternalServerError)
			return
		}

		// Cache content for 1 minute
		entry := CacheEntry{
			Content:     content,
			ContentType: resp.Header.Get("Content-Type"),
			Expiration:  time.Now().Add(1 * time.Minute),
		}
		cache.Set(url, entry)

		// Serve fetched content
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.Write(content)
		log.Println("Fetched and cached:", url)
	})

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
