package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


func ExploreLocations(name string) error {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", name)

	value, exists := cache.Get(url)
	if exists {

		var cacheResult locationArea
		err := json.Unmarshal(value, &cacheResult)
		if err != nil {
			return fmt.Errorf("error unmarshaling JSON: %w", err)
		}
		fmt.Printf("Exploring %v...\n", name)
		fmt.Println("Found Pokemon:")
		for _, result := range cacheResult.PokemonEncounters {
			fmt.Printf("- %v\n", result.Pokemon.Name)
		}
		return nil
	}

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

	var locationResult locationArea
	err = json.Unmarshal(body, &locationResult)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %w", err)
	}
	cache.Add(url, body)
	fmt.Printf("Exploring %v...\n", name)
	fmt.Println("Found Pokemon:")
	for _, result := range locationResult.PokemonEncounters {
		fmt.Printf("- %v\n", result.Pokemon.Name)
	}
	return nil
}
