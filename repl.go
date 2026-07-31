package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/paugomez86/gopokedex/internal/helpers"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
	args        []string
}

func startRepl(c *config) {
	scanner := bufio.NewScanner(os.Stdin)
	replCommands := getCommands()

	for {
		fmt.Print("Pokedex > ")
		for scanner.Scan() {
			input := helpers.CleanInput(scanner.Text())
			if command, ok := replCommands[input[0]]; !ok {
				fmt.Println("Unknown command")
			} else {
				var args []string
				for i := 1; i < len(input); i++ {
					args = append(args, input[i])
				}
				if err := command.callback(c, args); err != nil {
					fmt.Println(err)
				}
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
			args:        []string{"b"},
		},
		"explore": {
			name:        "explore",
			description: "Displays the list of pokemon likely to be found in the given area.",
			callback:    commandExplore,
			args:        []string{"area_name"},
		},
		"catch": {
			name:        "catch",
			description: "Trys to catch the given pokemon. It may fail!",
			callback:    commandCatch,
			args:        []string{"pokemon_name"},
		},
	}
}

// Quits the program.
func commandExit(c *config, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("Too many arguments")
	}
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// Displays the available commands.
func commandHelp(c *config, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("Too many arguments")
	}

	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, command := range getCommands() {
		fmt.Printf("%v", command.name)
		if command.args != nil {
			for _, arg := range command.args {
				fmt.Printf(" <%v>", arg)
			}
		}
		fmt.Printf(": %v\n", command.description)
	}
	fmt.Println()
	return nil
}

// Displays the 20 first location areas of the Pokemon world. Subsequent calls print the next 20 results.
func commandMap(c *config, args []string) error {
	// Argument handle
	if len(args) > 1 {
		return fmt.Errorf("Too many arguments")
	}
	back := false
	if slices.Contains(args, "b") {
		back = true
	}

	// The url handles the pagination along with the config struct.
	// config.pagination stores the next and/or previous urls to call in order to loop through the entire list.
	var url string
	if c.pagination.next == nil && c.pagination.previous == nil {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		if !back {
			if c.pagination.next != nil {
				url = *c.pagination.next
			} else {
				fmt.Println("you're on the last page")
				return nil
			}
		} else {
			if c.pagination.previous != nil {
				url = *c.pagination.previous
			} else {
				fmt.Println("you're on the first page")
				return nil
			}
		}
	}

	// Fetching resource
	var resource helpers.Resource = helpers.LocationArea{}

	data, err := resource.Unmarshal(url, c.cache)
	if err != nil {
		return err
	}
	locationArea, ok := data.(helpers.LocationArea)
	if !ok {
		return fmt.Errorf("Error decoding JSON response.\n")
	}

	c.pagination.next = locationArea.Next
	c.pagination.previous = locationArea.Previous
	areas := locationArea.Results

	// Displaying data
	for _, area := range areas {
		fmt.Printf("%v\n", area.Name)
	}
	return nil
}

// Displays the list of pokemon likely to be found in the given area in args[0].
func commandExplore(c *config, args []string) error {
	// Argument handle
	if len(args) != 1 {
		return fmt.Errorf("Expected 1 argument")
	}

	// Url handle
	url := "https://pokeapi.co/api/v2/location-area/" + args[0]

	// Fetching resource
	var resource helpers.Resource = helpers.PokemonEncounters{}

	data, err := resource.Unmarshal(url, c.cache)
	if err != nil {
		return err
	}
	pokemonEncounters, ok := data.(helpers.PokemonEncounters)
	if !ok {
		return fmt.Errorf("Error decoding JSON response.\n")
	}

	// Displaying data
	fmt.Printf("Exploring %v...\n", args[0])
	for _, p := range pokemonEncounters.PokemonEncounters {
		fmt.Printf(" - %v\n", p.Pokemon.Name)
	}

	return nil
}

func commandCatch(c *config, args []string) error {
	// Argument handle
	if len(args) != 1 {
		return fmt.Errorf("Expected 1 argument")
	}

	// Url handle
	url := "https://pokeapi.co/api/v2/pokemon/" + args[0]

	// Fetching resource
	var resource helpers.Resource = helpers.Pokemon{}

	val, err := resource.Unmarshal(url, c.cache)
	if err != nil {
		return err
	}
	pokemon, ok := val.(helpers.Pokemon)
	if !ok {
		return fmt.Errorf("Error decoding JSON response.\n")
	}

	// Trying to catch and display info
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)
	time.Sleep(time.Millisecond * 500)
	if helpers.TryCatchPokemon(pokemon) {
		fmt.Printf("%v was caught!\n", pokemon.Name)
		c.caught[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%v escaped!\n", pokemon.Name)
	}

	return nil
}
