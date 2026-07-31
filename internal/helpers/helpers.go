package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"

	"github.com/paugomez86/gopokedex/internal/pokecache"
)

type Resource interface {
	Unmarshal(string, *pokecache.Cache) (any, error)
}

type Pokemon struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}

type LocationAreas struct {
}

func (p Pokemon) Unmarshal(url string, c *pokecache.Cache) (any, error) {
	// Checking if url key is in cache
	var data []byte

	if cachedData, ok := c.Get(url); ok {
		data = cachedData
	} else {
		var err error
		if data, err = FetchResoruces(url); err != nil {
			return p, err
		}
		c.Add(url, data)
	}

	if err := json.Unmarshal(data, &p); err != nil {
		return p, fmt.Errorf("Error decoding JSON response: %v\n", err)
	}

	return p, nil
}

// Takes a string as input and returns a slice of its words using a whitespace as separator.
// The resulting words are lowercased and trimmed of leading and trailing whitespaces.
func CleanInput(input string) []string {
	var words []string
	for w := range strings.SplitSeq(input, " ") {
		if w != "" {
			words = append(words, strings.Trim(strings.ToLower(w), " "))
		}
	}
	return words
}

// Takes a string url and makes the GET request. Returns the raw data. Handles errors.
func FetchResoruces(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error fetching PokeApi: %v\n", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, fmt.Errorf("Resource not found")
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading JSON response: %v\n", err)
	}
	return data, nil
}

func TryCatchPokemon(p Pokemon) bool {
	if rand.Float64() >= (min(float64(p.BaseExperience), 300)/300)/2+0.1 {
		return true
	}
	return false
}
