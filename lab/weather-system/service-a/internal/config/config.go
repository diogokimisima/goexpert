package config

import (
	"os"
)

type Config struct {
	Port         string
	ServiceBURL  string
	OtelEndpoint string
	ServiceName  string
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serviceBURL := os.Getenv("SERVICE_B_URL")
	if serviceBURL == "" {
		serviceBURL = "http://service-b:8080"
	}

	otelEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otelEndpoint == "" {
		otelEndpoint = "otel-collector:4317"
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "service-a"
	}

	return &Config{
		Port:         port,
		ServiceBURL:  serviceBURL,
		OtelEndpoint: otelEndpoint,
		ServiceName:  serviceName,
	}, nil
}
