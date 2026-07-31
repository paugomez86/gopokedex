package main

import (
	"time"

	"github.com/paugomez86/gopokedex/internal/helpers"
	"github.com/paugomez86/gopokedex/internal/pokecache"
)

type config struct {
	pagination pagination
	cache      *pokecache.Cache
	caught     map[string]helpers.Pokemon
}

type pagination struct {
	next     *string
	previous *string
}

func main() {
	const cacheDuration = time.Millisecond * 15000

	c := config{
		pagination: pagination{
			next:     nil,
			previous: nil,
		},
		cache:  pokecache.NewCache(cacheDuration),
		caught: make(map[string]helpers.Pokemon),
	}
	startRepl(&c)
}
