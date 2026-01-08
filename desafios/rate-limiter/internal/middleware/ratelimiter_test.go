package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/config"
	"github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/limiter"
	"github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiterMiddleware_AllowsRequestsUnderLimit(t *testing.T) {
	store := storage.NewMemoryStorage()
	defer store.Close()

	cfg := &config.Config{
		IP:     config.RateLimitConfig{Requests: 5, Duration: time.Second, BlockDuration: 5 * time.Second},
		Token:  config.RateLimitConfig{Requests: 10, Duration: time.Second, BlockDuration: 5 * time.Second},
		Tokens: make(map[string]config.RateLimitConfig),
	}

	rl := limiter.New(store, cfg)
	middleware := RateLimiterMiddleware(rl)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	for i := 1; i <= 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Request %d should return 200", i)
		assert.Equal(t, "OK", w.Body.String(), "Request %d should return OK", i)
	}
}

func TestRateLimiterMiddleware_BlocksRequestsOverLimit(t *testing.T) {
	store := storage.NewMemoryStorage()
	defer store.Close()

	cfg := &config.Config{
		IP:     config.RateLimitConfig{Requests: 3, Duration: time.Second, BlockDuration: 5 * time.Second},
		Token:  config.RateLimitConfig{Requests: 10, Duration: time.Second, BlockDuration: 5 * time.Second},
		Tokens: make(map[string]config.RateLimitConfig),
	}

	rl := limiter.New(store, cfg)
	middleware := RateLimiterMiddleware(rl)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	for i := 1; i <= 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Request %d should return 200", i)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Request 4 should return 429")
	assert.Contains(t, w.Body.String(), MessageRateLimitExceeded)
}
