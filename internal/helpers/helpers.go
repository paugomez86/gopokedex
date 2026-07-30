package helpers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

type Pokemon struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
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
	if rand.Float64() >= (min(float64(p.BaseExperience))/300)/2+0.1 {
		return true
	}
	return false
}
