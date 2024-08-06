package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/briandowns/spinner"
	"impl/lib"
	"os"
	"time"
)

var (
	getCmd *flag.FlagSet
	setCmd *flag.FlagSet

	usageText = "Usage: get <city 1> <city 2> ... | set <key> <value>"
)

func init() {
	lib.EnsureConfig()

	getCmd = flag.NewFlagSet("get", flag.ExitOnError)
	setCmd = flag.NewFlagSet("set", flag.ExitOnError)
}

func main() {
	lib.HandleErr(errors.New("invalid usage"), len(os.Args) < 2)

	switch os.Args[1] {
	case "get":
		err := getCmd.Parse(os.Args[2:])
		lib.HandleErr(err, err != nil)

		handleGetCmd(getCmd.Args())
	case "set":
		err := setCmd.Parse(os.Args[2:])
		lib.HandleErr(err, err != nil)

		handleSetCmd(setCmd.Args())
	default:
		fmt.Println(usageText)
	}
}

func handleGetCmd(args []string) {
	length := len(args)
	lib.HandleErr(errors.New("invalid usage"), length == 0 || length > 5)

	apikey, err := lib.GetValue("apikey")
	lib.HandleErr(err, err != nil)
	units, err := lib.GetValue("units")
	lib.HandleErr(err, err != nil)

	getCityText := fmt.Sprintf("Getting city: %s", args[0])
	if length > 1 {
		cities := ""
		for _, val := range args {
			cities += " " + val
		}
		getCityText = fmt.Sprintf("Getting cities:%s", cities)
	}
	fmt.Println(getCityText)

	var data []lib.WeatherApi
	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	s.Start()
	for _, city := range args {
		url := fmt.Sprintf(
			"https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=%s",
			city,
			apikey,
			units,
		)
		fetched := lib.FetchWeatherData(url, s)
		data = append(data, fetched)
	}
	s.Stop()

	for _, content := range data {
		content.Print()
	}
}

func handleSetCmd(args []string) {
	lib.HandleErr(errors.New("invalid usage"), len(args) != 2)
	key, val := args[0], args[1]

	config := lib.UpdateConfig(key, val)
	configBytes, err := json.Marshal(config)
	lib.HandleErr(err, err != nil)

	err = os.WriteFile("config.json", configBytes, 0644)
	lib.HandleErr(err, err != nil)
}
