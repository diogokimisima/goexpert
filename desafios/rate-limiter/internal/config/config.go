package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Redis      RedisConfig
	IP         RateLimitConfig
	Token      RateLimitConfig
	Tokens     map[string]RateLimitConfig
	ServerPort string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type RateLimitConfig struct {
	Requests      int
	Duration      time.Duration
	BlockDuration time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		IP: RateLimitConfig{
			Requests:      getEnvAsInt("RATE_LIMIT_IP_REQUESTS", 5),
			Duration:      getEnvAsDuration("RATE_LIMIT_IP_DURATION", time.Second),
			BlockDuration: getEnvAsDuration("RATE_LIMIT_IP_BLOCK_DURATION", 5*time.Minute),
		},
		Token: RateLimitConfig{
			Requests:      getEnvAsInt("RATE_LIMIT_TOKEN_REQUESTS", 10),
			Duration:      getEnvAsDuration("RATE_LIMIT_TOKEN_DURATION", time.Second),
			BlockDuration: getEnvAsDuration("RATE_LIMIT_TOKEN_BLOCK_DURATION", 5*time.Minute),
		},
		ServerPort: getEnv("SERVER_PORT", "8080"),
		Tokens:     make(map[string]RateLimitConfig),
	}

	tokensConfig := getEnv("RATE_LIMIT_TOKENS", "")
	if tokensConfig != "" {
		tokens := strings.Split(tokensConfig, ",")
		for _, tokenStr := range tokens {
			parts := strings.Split(strings.TrimSpace(tokenStr), ":")
			if len(parts) == 4 {
				token := parts[0]
				requests, _ := strconv.Atoi(parts[1])
				duration, _ := time.ParseDuration(parts[2])
				blockDuration, _ := time.ParseDuration(parts[3])

				cfg.Tokens[token] = RateLimitConfig{
					Requests:      requests,
					Duration:      duration,
					BlockDuration: blockDuration,
				}
			}
		}
	}

	return cfg, nil
}

func (c *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
