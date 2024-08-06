package lib

import (
	"encoding/json"
	"fmt"
	"github.com/briandowns/spinner"
	"io"
	"net/http"
	"os"
)

type WeatherApi struct {
	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Name string `json:"name"`
}

type WeatherApiError struct {
	Message string `json:"message"`
}

func FetchWeatherData(url string, s *spinner.Spinner) WeatherApi {
	resp, err := http.Get(url)
	HandleErr(err, err != nil)

	body, err := io.ReadAll(resp.Body)
	HandleErr(err, err != nil)

	if resp.StatusCode > 299 {
		errData := WeatherApiError{}
		err := json.Unmarshal(body, &errData)
		HandleErr(err, err != nil)
		s.Stop()
		println(errData.Message)
		os.Exit(1)
	}

	err = resp.Body.Close()
	HandleErr(err, err != nil)

	apiData := WeatherApi{}
	err = json.Unmarshal(body, &apiData)
	HandleErr(err, err != nil)

	return apiData
}

func (data WeatherApi) Print() {
	fmt.Println("|")
	heading := fmt.Sprintf(
		"|%s - %s, %.0fC | feels like %.0fC",
		data.Name,
		data.Weather[0].Main,
		data.Main.Temp,
		data.Main.FeelsLike,
	)
	fmt.Println(heading)
	fmt.Println(fmt.Sprintf(
		"|Min temp: %.0f; Max temp: %.0f",
		data.Main.TempMin,
		data.Main.TempMax,
	))
	fmt.Println("|Wind speed:", data.Wind.Speed, "m/s")
	fmt.Println("|")
}
