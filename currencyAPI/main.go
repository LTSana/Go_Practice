package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/joho/godotenv"
)

type Rates struct {
	currency string
	rate     float64
}

type Currencies struct {
	Disclaimer string
	License    string
	Timestamp  string
	Base       string
	Rates      map[string]float64
	Rates2     []Rates
}

var (
	fromC  float64
	toC    float64
	amount string
)

func main() {

	// Load the environment file
	envFile, _ := godotenv.Read(".env")
	apiKey := envFile["API_KEY"]

	// Add the API key to the url
	url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s&base=USD", apiKey)

	// Make a request to the server
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// Set the headers
	r.Header.Add("accept", "application/json")

	// I don't know what this does
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer res.Body.Close() // Close the response once we done with it

	// Add the data to the structure
	var currenciesJson Currencies
	json.NewDecoder(res.Body).Decode(&currenciesJson)

	// Add the currency and rates
	for currency, rate := range currenciesJson.Rates {
		currenciesJson.Rates2 = append(currenciesJson.Rates2, Rates{currency: currency, rate: rate})
	}

	// Sort the options
	sort.Slice(currenciesJson.Rates2, func(i, j int) bool {
		return currenciesJson.Rates2[i].currency < currenciesJson.Rates2[j].currency
	})

	// Create the options for the form
	var huhOptions []huh.Option[float64]
	for _, o := range currenciesJson.Rates2 {
		huhOptions = append(huhOptions, huh.NewOption(o.currency, o.rate))
	}

	form := huh.NewForm(

		// For the currency
		huh.NewGroup(

			// Ask the user what currency they are converting from
			huh.NewSelect[float64]().
				Title("From Currency").
				Options(
					huhOptions...,
				).
				Value(&fromC). // Store the chosen option in the "fromC" variable
				Height(10),

			// Ask the user what currency they are converting to
			huh.NewSelect[float64]().
				Title("To Currency").
				Options(
					huhOptions...,
				).
				Value(&toC). // Store the chosen option in the "toC" variable
				Height(10),

			// Ask the user for the amount they want to convert
			huh.NewInput().
				Title("How much you want to convert?").
				CharLimit(99999).
				Placeholder("0.00").
				Value(&amount),
		),
	)

	// Run the form
	form.Run()

	// Parse the amount to a float 64 bit type for calculation
	_amount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Printf("Converted: %.2f\n", _amount*(toC/fromC))
}
