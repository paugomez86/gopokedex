package main

import (
	"time"

	"github.com/paugomez86/gopokedex/internal/helpers"
	"github.com/paugomez86/gopokedex/internal/pokecache"
)

type config struct {
	nextPage     *string
	previousPage *string
	cache        *pokecache.Cache
	caught       map[string]helpers.Pokemon
}

func main() {
	const cacheInterval = time.Millisecond * 5000

	c := config{
		nextPage:     nil,
		previousPage: nil,
		cache:        pokecache.NewCache(cacheInterval),
		caught:       make(map[string]helpers.Pokemon),
	}
	startRepl(&c)
}
