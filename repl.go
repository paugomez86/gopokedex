package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	next     string
	previous string
}

type locationArea struct {
	name string
	url  string
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	replCommands := getCommands()
	var c config

	for {
		fmt.Print("Pokedex > ")
		for scanner.Scan() {
			input := scanner.Text()
			if command, ok := replCommands[input]; !ok {
				fmt.Println("Unknown command")
			} else {
				command.callback(&c)
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
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays the available list of commands",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays 20 location areas. Subsequent map commands display the 20 next",
			callback:    commandMap,
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

// Displays 20 location areas of the Pokemon world.
func commandMap(c *config) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error fetching PokeApi: %v\n", err)
	}
	defer res.Body.Close()

	var data map[string]any
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("Error decoding JSON: %v\n", err)
	}

	if s, ok := data["next"].(string); ok {
		c.next = s
	}
	if s, ok := data["previous"].(string); ok {
		c.previous = s
	}
	areas := data["results"].([]locationArea)

	for _, area := range areas {
		fmt.Sprintf("%#v\n", area.name)
	}

	return nil
}
