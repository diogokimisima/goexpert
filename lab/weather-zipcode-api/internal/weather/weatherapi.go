package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type WeatherResponse struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

func GetWeather(apiKey string, city string) (float64, error) {
	if apiKey == "" {
		return 0, fmt.Errorf("weather API key is not set")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	baseURL := "https://api.weatherapi.com/v1/current.json"

	// Build the request URL
	params := url.Values{}
	params.Add("key", apiKey)
	params.Add("q", city)

	// Make the request
	resp, err := client.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		return 0, fmt.Errorf("error making request to Weather API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weather API returned status code %d", resp.StatusCode)
	}

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return 0, fmt.Errorf("error decoding Weather API response: %w", err)
	}

	return weather.Current.TempC, nil
}

func ConvertToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func ConvertToKelvin(celsius float64) float64 {
	return celsius + 273.15
}
