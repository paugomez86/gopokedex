package main

import (
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello    world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "123, testing",
			expected: []string{"123,", "testing"},
		},
		{
			input:    "HeLLo WorLd !!",
			expected: []string{"hello", "world", "!!"},
		},
	}
	for _, c := range cases {
		result := cleanInput(c.input)
		if len(result) != len(c.expected) {
			fmt.Printf("Expected length: %v - Actual: %v/n", len(c.expected), len(result))
		}
	}
}
