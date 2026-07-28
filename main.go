package main

import (
	"time"

	"github.com/paugomez86/gopokedex/pokecache"
)

type config struct {
	next     *string
	previous *string
	cache    *pokecache.Cache
}

func main() {
	c := config{
		next:     nil,
		previous: nil,
		cache:    pokecache.NewCache(time.Second * 10),
	}
	startRepl(&c)
}
