package config

import (
	"os"
)

type Config struct {
	Port          string
	ViaCepAPIKey  string
	WeatherAPIKey string
	WeatherAPIURL string
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:          port,
		ViaCepAPIKey:  os.Getenv("VIA_CEP_API_KEY"),
		WeatherAPIKey: os.Getenv("WEATHER_API_KEY"),
		WeatherAPIURL: os.Getenv("WEATHER_API_URL"),
	}, nil
}
