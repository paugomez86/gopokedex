package main

import (
	"time"

	"github.com/paugomez86/gopokedex/internal/pokecache"
)

type config struct {
	nextPage     *string
	previousPage *string
	cache        *pokecache.Cache
}

func main() {
	const cacheInterval = time.Millisecond * 5000

	c := config{
		nextPage:     nil,
		previousPage: nil,
		cache:        pokecache.NewCache(cacheInterval),
	}
	startRepl(&c)
}
