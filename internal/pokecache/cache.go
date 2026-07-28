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

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(t time.Duration) *Cache {
	c := Cache{
		entries:  make(map[string]cacheEntry),
		duration: t,
	}
	go c.reapLoop()
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if entry, ok := c.entries[key]; ok {
		return entry.val, true
	}
	return nil, false
}

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
