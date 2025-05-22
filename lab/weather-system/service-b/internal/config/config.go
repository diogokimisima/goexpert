package config

import (
	"os"
)

type Config struct {
	Port          string
	ViaCepAPIKey  string
	WeatherAPIKey string
	WeatherAPIURL string
	OtelEndpoint  string
	ServiceName   string
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	otelEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otelEndpoint == "" {
		otelEndpoint = "otel-collector:4317"
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "service-b"
	}

	return &Config{
		Port:          port,
		ViaCepAPIKey:  os.Getenv("VIA_CEP_API_KEY"),
		WeatherAPIKey: os.Getenv("WEATHER_API_KEY"),
		WeatherAPIURL: os.Getenv("WEATHER_API_URL"),
		OtelEndpoint:  otelEndpoint,
		ServiceName:   serviceName,
	}, nil
}
