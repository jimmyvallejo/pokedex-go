package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
)

var caughtPokemonMap = map[string]pokemonData{}

func ViewAllPokemon() error {
	fmt.Println("Your Pokedex:")
	for key, _ := range caughtPokemonMap {
		fmt.Printf("- %v\n", key)
	}
	return nil
}

func InspectPokemon(name string) error {
	if _, exists := caughtPokemonMap[name]; !exists {
		fmt.Printf("Pokemon %v not part of your Pokedex\n", name)
		return nil
	}
	foundPokemon := caughtPokemonMap[name]
	fmt.Printf("Name: %v\n", foundPokemon.Name)
	fmt.Printf("Height: %v\n", foundPokemon.Height)
	fmt.Printf("Weight: %v\n", foundPokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range foundPokemon.Stats {
		fmt.Printf("-%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pokeType := range foundPokemon.Types {
		fmt.Printf("-%v\n", pokeType.Type.Name)
	}
	return nil
}

func CatchPokemon(name string) error {

	if _, exists := caughtPokemonMap[name]; exists {
		fmt.Printf("Pokemon %v already caught\n", name)
		return nil
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v", name)

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making GET request: %w", err)
	}

	if res.StatusCode > 299 {
		return fmt.Errorf("HTTP request failed with status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return err
	}

	var pokemonResult pokemonData
	err = json.Unmarshal(body, &pokemonResult)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	randomNumber := rand.IntN(500)

	fmt.Printf("Throwing a ball at %v...\n", pokemonResult.Name)

	if randomNumber > pokemonResult.BaseExperience {
		caughtPokemonMap[name] = pokemonResult
		fmt.Printf("%v was caught! \n", pokemonResult.Name)
	} else {
		fmt.Printf("%v escaped! \n", pokemonResult.Name)
	}
	return nil
}
