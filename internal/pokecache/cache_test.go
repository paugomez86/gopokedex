package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const duration = 5 * time.Second

	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("test data"),
		},
		{
			key: "https://url.com/path",
			val: []byte("more data"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(duration)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("Expected to find key: %v\n", c.key)
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("Expected value: %v - Actual value: %v\n", c.val, val)
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 200 * time.Millisecond
	const waitTime = baseTime + 200*time.Millisecond

	c := struct {
		key string
		val []byte
	}{
		key: "https://example.com",
		val: []byte("test data"),
	}

	cache := NewCache(baseTime)
	cache.Add(c.key, c.val)

	if _, ok := cache.Get(c.key); !ok {
		t.Errorf("Expected to find key: %v\n", c.key)
		return
	}

	time.Sleep(waitTime)

	if _, ok := cache.Get(c.key); ok {
		t.Errorf("Expected not to find key: %v", c.key)
	}
}
