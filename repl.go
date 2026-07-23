package main

import "strings"

// Takes a string as input and returns a slice of its words using a whitespace as separator.
// The resulting words are lowercased and trimmed of leading and trailing whitespaces.
func cleanInput(input string) []string {
	var words []string
	for w := range strings.SplitSeq(input, " ") {
		if w != "" {
			words = append(words, strings.Trim(strings.ToLower(w), " "))
		}
	}
	return words
}
