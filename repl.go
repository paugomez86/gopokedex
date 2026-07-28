package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func startRepl(c *config) {
	scanner := bufio.NewScanner(os.Stdin)
	replCommands := getCommands()

	for {
		fmt.Print("Pokedex > ")
		for scanner.Scan() {
			input := scanner.Text()
			if command, ok := replCommands[input]; !ok {
				fmt.Println("Unknown command")
			} else {
				command.callback(c)
			}

			if err := scanner.Err(); err == nil {
				break
			}
		}
	}
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex.",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays the available list of commands.",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays 20 location areas. Subsequent map commands display the 20 next.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations when used after the map command.",
			callback:    commandMapb,
		},
	}
}

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

// Quits the program.
func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// Displays the available commands.
func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, v := range getCommands() {
		fmt.Printf("%v: %v\n", v.name, v.description)
	}
	fmt.Println()
	return nil
}

// Displays the 20 first location areas of the Pokemon world. Subsequent calls print the next 20 results.
func commandMap(c *config) error {
	type Response struct {
		Next     *string `json:"next"`
		Previous *string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
		} `json:"results"`
	}

	// The url handles the pagination along with the config struct.
	// config stores the next and/or previous urls to call in order to loop through the entire list.
	var url string
	if c.next == nil && c.previous == nil {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		if c.next != nil {
			url = *c.next
		} else {
			fmt.Println("you're on the last page")
			return nil
		}
	}

	// Cache checking
	var response Response
	var data []byte

	if val, ok := c.cache.Get(url); ok {
		data = val
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error fetching PokeApi: %v\n", err)
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading JSON response: %v\n", err)
		}

		c.cache.Add(url, data)
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("Error decodiing JSON response: %v\n", err)
	}

	c.next = response.Next
	c.previous = response.Previous
	areas := response.Results

	for _, area := range areas {
		fmt.Println(area.Name)
	}
	return nil
}

// Displays the previous 20 locations if used after commandMap.
func commandMapb(c *config) error {
	type Response struct {
		Next     *string `json:"next"`
		Previous *string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
		} `json:"results"`
	}

	var url string
	if c.previous != nil {
		url = *c.previous
	} else {
		fmt.Println("you're on the first page")
		return nil
	}

	// Cache checking
	var response Response
	var data []byte

	if val, ok := c.cache.Get(url); ok {
		data = val
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error fetching PokeApi: %v\n", err)
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading JSON response: %v\n", err)
		}

		c.cache.Add(url, data)
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("Error decodiing JSON response: %v\n", err)
	}

	c.next = response.Next
	c.previous = response.Previous
	areas := response.Results

	for _, area := range areas {
		fmt.Println(area.Name)
	}
	return nil
}
