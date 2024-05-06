package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

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
}

// NewCache creates a new instance of Cache
func NewCache() *Cache {
	return &Cache{
		entries: make(map[string]CacheEntry),
	}
}

// Get retrieves a cached entry by key
func (c *Cache) Get(key string) (CacheEntry, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, ok := c.entries[key]
	return entry, ok
}

// Set inserts or updates a cached entry
func (c *Cache) Set(key string, entry CacheEntry) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.entries[key] = entry
}

func main() {
	
	cache := NewCache()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
			return
		}

		if entry, ok := cache.Get(url); ok && time.Now().Before(entry.Expiration) {
			// Serve cached content
			w.Header().Set("Content-Type", entry.ContentType)
			w.Write(entry.Content)
			fmt.Println("Served from cache:", url)
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
		fmt.Println("Fetched and cached:", url)
	})

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
