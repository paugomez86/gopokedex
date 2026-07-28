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
	return &c
}

func (c *Cache) Add(key string, val []byte) error {
	return nil
}
