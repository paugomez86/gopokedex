package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/paugomez86/gopokedex/internal/helpers"
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
			expected: []string(nil),
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

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			actual := helpers.CleanInput(c.input)
			expected := c.expected
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("Expected value: %#v, actual value: %#v", c.expected, actual)
			}
		})
	}
}
