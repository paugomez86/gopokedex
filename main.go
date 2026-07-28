package main

import (
	"time"

	"github.com/paugomez86/gopokedex/internal/pokecache"
)

type config struct {
	next     *string
	previous *string
	cache    *pokecache.Cache
}

func main() {
	const cacheInterval = time.Millisecond * 5000

	c := config{
		next:     nil,
		previous: nil,
		cache:    pokecache.NewCache(cacheInterval),
	}
	startRepl(&c)
}
