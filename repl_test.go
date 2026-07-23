package main

import (
	"reflect"
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
		actual := cleanInput(c.input)
		expected := c.expected
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Test %v: expected %#v, actual %#v", i+1, c.expected, actual)
		}
	}
}
