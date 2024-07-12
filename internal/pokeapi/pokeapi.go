package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/jimmyvallejo/pokedex-go/internal/pokecache"
)

type config struct {
	next     *string
	previous *string
}

var locationConfig = config{next: nil, previous: nil}

type responseResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type locationResponse struct {
	Count    int              `json:"count"`
	Next     *string          `json:"next"`
	Previous *string          `json:"previous"`
	Results  []responseResult `json:"results"`
}

var cache *pokecache.Cache

func InitWithCache(c *pokecache.Cache) {
	cache = c
}

func FetchLocations() error {
	url := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"

	if locationConfig.next != nil {
		url = *locationConfig.next
	}

	value, exists := cache.Get(url)
	if exists {

		var cacheResult locationResponse
		err := json.Unmarshal(value, &cacheResult)
		if err != nil {
			return fmt.Errorf("error unmarshaling JSON: %w", err)
		}
		fmt.Println("cache used")
		locationConfig.previous = cacheResult.Previous
		locationConfig.next = cacheResult.Next
		for _, result := range cacheResult.Results {
			fmt.Println(result.Name)
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

	var locationResult locationResponse

	err = json.Unmarshal(body, &locationResult)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	locationConfig.previous = locationResult.Previous
	locationConfig.next = locationResult.Next
	cache.Add(url, body)

	for _, result := range locationResult.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func FetchPrevious() error {
	if locationConfig.previous == nil {
		return errors.New("no previous results available")
	}
	url := *locationConfig.previous

	value, exists := cache.Get(url)
	if exists {

		var cacheResult locationResponse
		err := json.Unmarshal(value, &cacheResult)
		if err != nil {
			return fmt.Errorf("error unmarshaling JSON: %w", err)
		}
		fmt.Println("cache used")
		locationConfig.previous = cacheResult.Previous
		locationConfig.next = cacheResult.Next
		for _, result := range cacheResult.Results {
			fmt.Println(result.Name)
		}
		return nil
	}

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making GET request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return fmt.Errorf("HTTP request failed with status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	var locationResult locationResponse
	err = json.Unmarshal(body, &locationResult)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	locationConfig.previous = locationResult.Previous
	locationConfig.next = locationResult.Next

	for _, result := range locationResult.Results {
		fmt.Println(result.Name)
	}

	return nil
}
