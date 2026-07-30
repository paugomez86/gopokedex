package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
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
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations when used after the map command.",
			callback:    commandMapb,
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
	type Response struct {
		Next     *string `json:"next"`
		Previous *string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
		} `json:"results"`
	}

	if len(args) > 0 {
		return fmt.Errorf("Too many arguments")
	}

	// The url handles the pagination along with the config struct.
	// config stores the next and/or previous urls to call in order to loop through the entire list.
	var url string
	if c.pagination.next == nil && c.pagination.previous == nil {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		if c.pagination.next != nil {
			url = *c.pagination.next
		} else {
			fmt.Println("you're on the last page")
			return nil
		}
	}

	// Checking if url key is in cache
	var data []byte

	if cachedData, ok := c.cache.Get(url); ok {
		data = cachedData
	} else {
		var err error
		if data, err = helpers.FetchResoruces(url); err != nil {
			return err
		}
		c.cache.Add(url, data)
	}

	var response Response
	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("Error decoding JSON response: %v\n", err)
	}

	c.pagination.next = response.Next
	c.pagination.previous = response.Previous
	areas := response.Results

	for _, area := range areas {
		fmt.Printf("%v\n", area.Name)
	}
	return nil
}

// Displays the previous 20 locations if used after commandMap.
func commandMapb(c *config, args []string) error {
	type Response struct {
		Next     *string `json:"next"`
		Previous *string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
		} `json:"results"`
	}

	if len(args) > 0 {
		return fmt.Errorf("Too many arguments")
	}

	var url string
	if c.pagination.next != nil {
		url = *c.pagination.previous
	} else {
		fmt.Println("You're on the first page")
		return nil
	}

	// Checking if url key is in cache
	var data []byte

	if cachedData, ok := c.cache.Get(url); ok {
		data = cachedData
	} else {
		var err error
		if data, err = helpers.FetchResoruces(url); err != nil {
			return err
		}
		c.cache.Add(url, data)
	}

	var response Response
	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("Error decoding JSON response: %v\n", err)
	}

	c.pagination.next = response.Next
	c.pagination.previous = response.Previous
	areas := response.Results

	for _, area := range areas {
		fmt.Printf("%v\n", area.Name)
	}
	return nil
}

// Displays the list of pokemon likely to be found in the given area in args[0].
func commandExplore(c *config, args []string) error {
	type Response struct {
		PokemonEncounters []struct {
			Pokemon struct {
				Name string `json:"name"`
			} `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}

	if len(args) != 1 {
		return fmt.Errorf("Expected 1 argument")
	}

	url := "https://pokeapi.co/api/v2/location-area/" + args[0]

	var data []byte

	if cachedData, ok := c.cache.Get(url); ok {
		data = cachedData
	} else {
		var err error
		if data, err = helpers.FetchResoruces(url); err != nil {
			return err
		}
		c.cache.Add(url, data)
	}

	var response Response
	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("Error decoding JSON response: %v\n", err)
	}

	pokemon := response.PokemonEncounters

	fmt.Printf("Exploring %v...\n", args[0])
	for _, p := range pokemon {
		fmt.Printf(" - %v\n", p.Pokemon.Name)
	}

	return nil
}

func commandCatch(c *config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected 1 argument")
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + args[0]

	// Checking if url key is in cache
	/* var data []byte

	if cachedData, ok := c.cache.Get(url); ok {
		data = cachedData
	} else {
		var err error
		if data, err = helpers.FetchResoruces(url); err != nil {
			return err
		}
		c.cache.Add(url, data)
	}

	var pokemon helpers.Pokemon
	if err := json.Unmarshal(data, &pokemon); err != nil {
		return fmt.Errorf("Error decoding JSON response: %v\n", err)
	}
	*/

	var resource helpers.Resource = helpers.Pokemon{}

	val, err := resource.Unmarshal(url, c.cache)
	if err != nil {
		return err
	}
	pokemon, ok := val.(helpers.Pokemon)
	if !ok {
		return fmt.Errorf("Error decoding JSON response.\n")
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)
	time.Sleep(time.Second * 1)
	if helpers.TryCatchPokemon(pokemon) {
		fmt.Printf("%v was caught!\n", pokemon.Name)
		c.caught[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%v escaped!\n", pokemon.Name)
	}

	return nil
}
