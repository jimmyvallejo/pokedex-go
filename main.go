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

type Commander interface {
	Execute(args []string) error
}

type cliCommand struct {
	name        string
	description string
	commander   Commander
}
type SimpleCommander struct {
	callback func() error
}

func (sc SimpleCommander) Execute(args []string) error {
	return sc.callback()
}

type ArgCommander struct {
	callback func(string) error
}

func (ac ArgCommander) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("this command requires an argument")
	}
	return ac.callback(args[0])
}

var commandMap map[string]cliCommand

func init() {
	commandMap = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message.",
			commander:   SimpleCommander{callback: commandHelp},
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex.",
			commander:   SimpleCommander{callback: commandExit},
		},
		"map": {
			name:        "map",
			description: "Displays the name of 20 locations in the Pokemon World.",
			commander:   SimpleCommander{callback: pokeapi.FetchLocations},
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 locations in the Pokemon World.",
			commander:   SimpleCommander{callback: pokeapi.FetchPrevious},
		},
		"explore": {
			name:        "explore",
			description: "Explore the different Pokemon living in a specific location.",
			commander:   ArgCommander{callback: pokeapi.ExploreLocations},
		},
		"catch": {
			name:        "catch",
			description: "Type the name of a Pokemon to try to catch it and add to your Pokedex.",
			commander:   ArgCommander{callback: pokeapi.CatchPokemon},
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon in your Pokedex to view its data.",
			commander:   ArgCommander{callback: pokeapi.InspectPokemon},
		},
		"pokedex": {
			name:        "pokedex",
			description: "View all Pokemon you have caught.",
			commander:   SimpleCommander{callback: pokeapi.ViewAllPokemon},
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

		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}

		commandName := args[0]

		if cmd, exists := commandMap[commandName]; exists {
			err := cmd.commander.Execute(args[1:])
			if err != nil {
				fmt.Println("An error occurred handling your request:", err)
			}
		} else {
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}
