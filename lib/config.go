package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Key string `json:"key"`
	//Theme string `json:"theme"`
	Units string `json:"units"`
}

var UNITS = []string{"standard", "metric", "imperial"}

var config = Config{}

func EnsureConfig() {
	_, err := os.Stat("config.json")
	if os.IsNotExist(err) {
		file, err := os.Create("config.json")
		HandleErr(err, err != nil)

		defer func() {
			err := file.Close()
			HandleErr(err, err != nil)
		}()

		_, err = file.Write([]byte("{\"units\":\"metric\"}"))
		HandleErr(err, err != nil)
	}
}

func getConfig() Config {
	content, err := os.ReadFile("config.json")
	HandleErr(err, err != nil)

	err = json.Unmarshal(content, &config)
	HandleErr(err, err != nil)

	return config
}

func UpdateConfig(key string, val string) Config {
	getConfig()

	switch key {
	case "apikey":
		config.Key = val
	default:
		fmt.Println("unknown key")
		os.Exit(1)
	}

	return config
}

func GetValue(key string) (string, error) {
	getConfig()

	switch key {
	case "apikey":
		if config.Key == "" {
			return "", errors.New("API key is required.\nRun: set apikey <key>")
		}
		return config.Key, nil
	case "units":
		if config.Units == "" {
			return "", errors.New("units option is required")
		}
		if !Contains(UNITS, config.Units) {
			return "", errors.New(
				"Units option is invalid. Expected one of this options: " +
					Join(UNITS, " | ") +
					"; got " +
					config.Units +
					"\nRun: set units <option>",
			)
		}
		return config.Units, nil
	default:
		return "", errors.New("unknown key")
	}
}
