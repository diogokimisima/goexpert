package main

import (
    "log"
    "net/http"
    "weather-zipcode-api/internal/api"
    "weather-zipcode-api/internal/config"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("could not load config: %v", err)
    }

    router := api.NewRouter(cfg)

    log.Printf("Starting server on port %s...", cfg.Port)
    if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
        log.Fatalf("could not start server: %v", err)
    }
}