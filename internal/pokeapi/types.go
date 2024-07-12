package pokeapi

// Map Location Structs

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

// Explore Area structs

type locationArea struct {
	PokemonEncounters []pokemonEncounter `json:"pokemon_encounters"`
}

type pokemonEncounter struct {
	Pokemon pokemon `json:"pokemon"`
}

type pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Catch/Inspect Pokemon

type pokemonData struct {
    ID             int    `json:"id"`
    Name           string `json:"name"`
    BaseExperience int    `json:"base_experience"`
    Height         int    `json:"height"`
    Weight         int    `json:"weight"`
    Stats []struct {
        BaseStat int `json:"base_stat"`
        Effort   int `json:"effort"`
        Stat     struct {
            Name string `json:"name"`
            URL  string `json:"url"`
        } `json:"stat"`
    } `json:"stats"`
    Types []struct {
        Slot int `json:"slot"`
        Type struct {
            Name string `json:"name"`
            URL  string `json:"url"`
        } `json:"type"`
    } `json:"types"`
}