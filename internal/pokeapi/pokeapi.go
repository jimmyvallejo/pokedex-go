package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

func FetchLocations() error {
	url := "https://pokeapi.co/api/v2/location-area"

	if locationConfig.next != nil {
		url = *locationConfig.next
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
