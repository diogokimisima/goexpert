package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/config"
	"github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/limiter"
	"github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/middleware"
	"github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/storage"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Carrega as configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Inicializa o storage (Redis)
	store, err := storage.NewRedisStorage(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer store.Close()

	log.Println("Connected to Redis successfully")

	// Cria o rate limiter
	rateLimiter := limiter.New(store, cfg)

	// Configura o roteador
	r := chi.NewRouter()

	// Middlewares globais
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)

	// Aplica o middleware de rate limiting
	r.Use(middleware.RateLimiterMiddleware(rateLimiter))

	// Rotas
	r.Get("/", handleHome)
	r.Get("/health", handleHealth)
	r.Post("/api/data", handleData)
	r.Get("/api/info", handleInfo)

	// Inicia o servidor
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Printf("Rate Limit - IP: %d req/%v, Block: %v",
		cfg.IP.Requests, cfg.IP.Duration, cfg.IP.BlockDuration)
	log.Printf("Rate Limit - Token (default): %d req/%v, Block: %v",
		cfg.Token.Requests, cfg.Token.Duration, cfg.Token.BlockDuration)

	if len(cfg.Tokens) > 0 {
		log.Printf("Custom token limits configured for %d tokens", len(cfg.Tokens))
		for token, limit := range cfg.Tokens {
			log.Printf("  - Token %s: %d req/%v, Block: %v",
				token, limit.Requests, limit.Duration, limit.BlockDuration)
		}
	}

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Welcome to Rate Limiter API",
		"endpoints": []string{
			"GET /health - Health check",
			"GET /api/info - Get API information",
			"POST /api/data - Submit data",
		},
		"usage": map[string]string{
			"rate_limit_by_ip":    "Requests are limited by IP address",
			"rate_limit_by_token": "Use API_KEY header to authenticate and get higher limits",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleData(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"message": "Data received successfully",
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("API_KEY")

	response := map[string]interface{}{
		"message":       "API Information",
		"version":       "1.0.0",
		"authenticated": apiKey != "",
	}

	if apiKey != "" {
		response["token"] = apiKey
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
