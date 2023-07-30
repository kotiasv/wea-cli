package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/fatih/color"
)

// Weather data format
type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	// sub commands
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	setCmd := flag.NewFlagSet("set", flag.ExitOnError)

	getLocation := getCmd.String("loc", "", "Name of the location (Moscow is default)")
	setKey := setCmd.String("key", "", "Api key")

	checkFile(file)

	// check for args
	if len(os.Args) < 2 {
		color.Red("To use the CLI you need to provide api key")
		message := fmt.Sprint(
			"Wea - weather app inside your console\n\n",
			"Commands\n",
			"  get: shows weather in the default location (by default its Elektrostal)\n",
			"  get -loc <LOCATION>: shows weather in the provided location\n\n",
			"  set -key <KEY>: set api key",
		)
		color.Magenta(message)
		return
	}

	// track subcommands
	switch os.Args[1] {
	case "get":
		HandleGet(getCmd, getLocation)
	case "set":
		HandleSet(setCmd, setKey)
	default:
		color.Red("expected 'get' or 'set' subcommands")
		return
	}
}

func HandleGet(getCmd *flag.FlagSet, loc *string) {
	getCmd.Parse(os.Args[2:])

	apiKey := getKey()
	if len(apiKey) == 0 {
		color.Red("Api key is not provided")
		color.Magenta("You can get your key here: https://www.weatherapi.com/\n\nTo set the key, run the command:")
		color.White("wea set -key <KEY>")

		return
	}
	apiUrl := "http://api.weatherapi.com/v1/forecast.json?key=" + apiKey + "&days=1&aqi=no&alerts=no"

	q := "Elektrostal"
	if *loc != "" {
		q = *loc
	}

	res, err := http.Get(apiUrl + "&q=" + q)
	check(err)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		message := fmt.Sprint(
			"Weather API is not available or Api Url is invalid\n",
			"Check if your Api Key is correct\n\n",
			"Current value:",
		)
		color.Red(message)
		color.White(apiKey)
		return
	}

	body, err := io.ReadAll(res.Body)
	check(err)

	var weather Weather
	err = json.Unmarshal(body, &weather)
	check(err)

	showWeather(weather)
}

func HandleSet(setCmd *flag.FlagSet, key *string) {
	setCmd.Parse(os.Args[2:])
	if *key == "" {
		color.Red("key must be provided")
		return
	} else {
		color.Magenta("Api key set as:")
		color.White(*key)
	}
	setKey(*key)
}
