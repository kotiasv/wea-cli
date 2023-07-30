package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Data struct {
	Key string `json:"key"`
}

func check(err interface{}) {
	if err != nil {
		panic(err)
	}
}

var ex, _ = os.Executable()
var file = strings.Replace(ex, "wea.exe", "key.json", 1)

// Create json file if it's not found
func checkFile(filePath string) {
	_, err := os.ReadFile(file)
	if err != nil {
		_, err := os.Create("key.json")
		check(err)

		dummyData := Data{""}
		res, err := json.Marshal(dummyData)
		check(err)

		err = os.WriteFile(file, res, 0644)
		check(err)
	}
}

// Methods
func getKey() string {
	res, err := os.ReadFile(file)
	check(err)

	var data Data
	err = json.Unmarshal(res, &data)
	check(err)

	return data.Key
}

func setKey(key string) {
	data := Data{key}
	res, err := json.Marshal(data)
	check(err)

	err = os.WriteFile(file, res, 0644)
	check(err)
}

// App functionality
func showWeather(weather Weather) {
	location, current, hours :=
		weather.Location,
		weather.Current,
		weather.Forecast.Forecastday[0].Hour

	fmt.Printf(
		"%s, %s: %.0fC, %s\n",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf(
			"%s - %.0fC, %.0f%%, %s\n",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		fmt.Print(message)
	}
}
