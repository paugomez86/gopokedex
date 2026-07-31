package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu       sync.Mutex
	entries  map[string]cacheEntry
	duration time.Duration
}

// Each entry has the resource URL and the data returned bu PokeAPI
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Initializes the catche struct, runs reapLoop and returns the pointer
func NewCache(t time.Duration) *Cache {
	c := Cache{
		entries:  make(map[string]cacheEntry),
		duration: t,
	}
	go c.reapLoop()
	return &c
}

// Adds an entry struct to the cache
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.entries[key] = entry
}

// Gets the url and checks if it's already cached
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if entry, ok := c.entries[key]; ok {
		return entry.val, true
	}
	return nil, false
}

// Redundant function that checks entries every c.duration time and removes them if they are old enough
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.duration)
	for range ticker.C {
		for k, entry := range c.entries {
			c.mu.Lock()
			if time.Since(entry.createdAt) >= c.duration {
				delete(c.entries, k)
			}
			c.mu.Unlock()
		}
	}
}
