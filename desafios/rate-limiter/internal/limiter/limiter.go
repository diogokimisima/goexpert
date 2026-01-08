package limiter

import (
	"context"
	"fmt"

	"github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/config"
	"github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/storage"
)

type RateLimiter struct {
	storage storage.Storage
	config  *config.Config
}

func New(storage storage.Storage, cfg *config.Config) *RateLimiter {
	return &RateLimiter{
		storage: storage,
		config:  cfg,
	}
}

func (rl *RateLimiter) CheckLimit(ctx context.Context, identifier, token string) (bool, error) {
	var key string
	var limit config.RateLimitConfig

	if token != "" {
		key = fmt.Sprintf("token:%s", token)

		if tokenConfig, exists := rl.config.Tokens[token]; exists {
			limit = tokenConfig
		} else {
			limit = rl.config.Token
		}
	} else {
		key = fmt.Sprintf("ip:%s", identifier)
		limit = rl.config.IP
	}

	blocked, err := rl.storage.IsBlocked(ctx, key)
	if err != nil {
		return false, fmt.Errorf("error checking block status: %w", err)
	}
	if blocked {
		return false, nil
	}

	count, err := rl.storage.Increment(ctx, key, limit.Duration)
	if err != nil {
		return false, fmt.Errorf("error incrementing counter: %w", err)
	}

	if count > int64(limit.Requests) {
		if err := rl.storage.SetBlock(ctx, key, limit.BlockDuration); err != nil {
			return false, fmt.Errorf("error setting block: %w", err)
		}
		return false, nil
	}

	return true, nil
}

func (rl *RateLimiter) GetRemainingRequests(ctx context.Context, identifier, token string) (int64, error) {
	var key string
	var limit config.RateLimitConfig

	if token != "" {
		key = fmt.Sprintf("token:%s", token)
		if tokenConfig, exists := rl.config.Tokens[token]; exists {
			limit = tokenConfig
		} else {
			limit = rl.config.Token
		}
	} else {
		key = fmt.Sprintf("ip:%s", identifier)
		limit = rl.config.IP
	}

	count, err := rl.storage.Get(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("error getting counter: %w", err)
	}

	remaining := int64(limit.Requests) - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}
