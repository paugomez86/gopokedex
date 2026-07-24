package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	replCommands := getCommands()

	for {
		fmt.Print("Pokedex > ")
		for scanner.Scan() {
			input := scanner.Text()
			if command, ok := replCommands[input]; !ok {
				fmt.Println("Unknown command")
			} else {
				command.callback()
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

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, v := range getCommands() {
		fmt.Printf("%v: %v\n", v.name, v.description)
	}
	fmt.Println()
	return nil
}
