package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jimmyvallejo/pokedex-go/internal/pokeapi"
	"github.com/jimmyvallejo/pokedex-go/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commandMap map[string]cliCommand

func init() {
	commandMap = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the name of 20 locations in the Pokemon World",
			callback:    pokeapi.FetchLocations,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 locations in the Pokemon World",
			callback:    pokeapi.FetchPrevious,
		},
	}
}

func commandExit() error {
	fmt.Println("Exiting the Pokedex")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range commandMap {
		_, err := fmt.Printf("%v: %v\n", cmd.name, cmd.description)
		if err != nil {
			return fmt.Errorf("error printing command info: %w", err)
		}
	}
	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	cache := pokecache.InitCache(60 * time.Second)
	pokeapi.InitWithCache(cache)

	for {
		fmt.Print("Pokedex > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input")
			continue
		}
		input = strings.TrimSuffix(input, "\n")
		input = strings.ToLower(input)

		if cmd, exists := commandMap[input]; exists {
			err = cmd.callback()
			if err != nil {
				fmt.Println("An error occurred handling your request:", err)
			}
		} else {
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}
