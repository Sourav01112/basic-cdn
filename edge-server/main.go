package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type CacheItem struct {
	Data      []byte
	ExpiresAt time.Time
}

type Cache struct {
	data map[string]*CacheItem
	mu   sync.RWMutex
	ttl  time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	cache := &Cache{
		data: make(map[string]*CacheItem),
		ttl:  ttl,
	}

	go cache.cleanup()

	return cache
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		log.Printf("Cache: Item %s expired at %s", key, item.ExpiresAt.Format(time.RFC3339))
		return nil, false
	}

	return item.Data, true
}

func (c *Cache) Set(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiresAt := time.Now().Add(c.ttl)

	c.data[key] = &CacheItem{
		Data:      data,
		ExpiresAt: expiresAt, // type: time.Time
	}

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		loc = time.FixedZone("IST", 5*3600+30*60)
	}
	log.Printf("Cache: Stored %s, expires at %s",
		key, expiresAt.In(loc).Format(time.RFC3339))
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(30 * time.Second) // every 30 seconds
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		expired := []string{}

		for key, item := range c.data {
			if now.After(item.ExpiresAt) {
				expired = append(expired, key)
			}
		}

		for _, key := range expired {
			delete(c.data, key)
			log.Printf("Cache: Cleaned up expired item %s", key)
		}

		if len(expired) > 0 {
			log.Printf("Cache: Cleaned up %d expired items", len(expired))
		}

		c.mu.Unlock()
	}
}

func (c *Cache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := len(c.data)
	active := 0
	expired := 0
	now := time.Now()

	for _, item := range c.data {
		if now.Before(item.ExpiresAt) {
			active++
		} else {
			expired++
		}
	}

	return map[string]interface{}{
		"total_items":   total,
		"active_items":  active,
		"expired_items": expired,
		"ttl_seconds":   int(c.ttl.Seconds()),
	}
}

var cache *Cache

func main() {
	originURL := os.Getenv("ORIGIN_URL")
	if originURL == "" {
		log.Fatal("ORIGIN_URL environment variable required")
	}

	cacheTTL := 60 * time.Second
	cache = NewCache(cacheTTL)

	log.Printf("Edge Server: Cache TTL set to %v", cacheTTL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if cachedData, exists := cache.Get(path); exists {
			log.Printf("Edge Server: Cache HIT for %s", path)

			var cachedJSON map[string]interface{}
			if err := json.Unmarshal(cachedData, &cachedJSON); err == nil {
				cachedJSON["cached"] = true
				cachedJSON["cache_server"] = "edge-server"

				modifiedData, _ := json.Marshal(cachedJSON)

				w.Header().Set("X-Cache", "HIT")
				w.Header().Set("Content-Type", "application/json")
				w.Write(modifiedData)
			} else {
				w.Header().Set("X-Cache", "HIT")
				w.Write(cachedData)
			}
			return
		}

		log.Printf("Edge Server: Cache MISS for %s, fetching from origin", path)

		resp, err := http.Get(originURL + path)
		if err != nil {
			http.Error(w, "Origin server error", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Read error", http.StatusInternalServerError)
			return
		}

		cache.Set(path, data)

		w.Header().Set("X-Cache", "MISS")
		w.Header().Set("X-Cache-TTL", cacheTTL.String())
		w.Header().Set("X-CDN-Server", "edge-server")
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.Write(data)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Edge Server OK"))
	})

	log.Println("Edge Server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
