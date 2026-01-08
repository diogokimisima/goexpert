package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/config"
	"github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateLimiter_IPBasedLimiting(t *testing.T) {
	store := storage.NewMemoryStorage()
	defer store.Close()

	cfg := &config.Config{
		IP: config.RateLimitConfig{
			Requests:      5,
			Duration:      time.Second,
			BlockDuration: 5 * time.Second,
		},
		Token:  config.RateLimitConfig{Requests: 10, Duration: time.Second, BlockDuration: 5 * time.Second},
		Tokens: make(map[string]config.RateLimitConfig),
	}

	rl := New(store, cfg)
	ctx := context.Background()
	ip := "192.168.1.1"

	for i := 1; i <= 5; i++ {
		allowed, err := rl.CheckLimit(ctx, ip, "")
		require.NoError(t, err)
		assert.True(t, allowed, "Request %d should be allowed", i)
	}

	allowed, err := rl.CheckLimit(ctx, ip, "")
	require.NoError(t, err)
	assert.False(t, allowed, "Request 6 should be blocked")
}

func TestRateLimiter_TokenBasedLimiting(t *testing.T) {
	store := storage.NewMemoryStorage()
	defer store.Close()

	cfg := &config.Config{
		IP:     config.RateLimitConfig{Requests: 5, Duration: time.Second, BlockDuration: 5 * time.Second},
		Token:  config.RateLimitConfig{Requests: 10, Duration: time.Second, BlockDuration: 5 * time.Second},
		Tokens: make(map[string]config.RateLimitConfig),
	}

	rl := New(store, cfg)
	ctx := context.Background()
	ip := "192.168.1.1"
	token := "test-token"

	for i := 1; i <= 10; i++ {
		allowed, err := rl.CheckLimit(ctx, ip, token)
		require.NoError(t, err)
		assert.True(t, allowed, "Request %d should be allowed", i)
	}

	allowed, err := rl.CheckLimit(ctx, ip, token)
	require.NoError(t, err)
	assert.False(t, allowed, "Request 11 should be blocked")
}
