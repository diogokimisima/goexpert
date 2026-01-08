package middleware

import (
	"net/http"
	"strings"

	"github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/limiter"
)

const (
	HeaderAPIKey             = "API_KEY"
	MessageRateLimitExceeded = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

func RateLimiterMiddleware(rl *limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)
			token := r.Header.Get(HeaderAPIKey)

			allowed, err := rl.CheckLimit(r.Context(), ip, token)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"message":"` + MessageRateLimitExceeded + `"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}

	return ip
}
