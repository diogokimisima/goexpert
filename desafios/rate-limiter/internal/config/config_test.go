package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_WithDefaults(t *testing.T) {
	os.Clearenv()
	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "localhost", cfg.Redis.Host)
	assert.Equal(t, "6379", cfg.Redis.Port)
	assert.Equal(t, 5, cfg.IP.Requests)
	assert.Equal(t, time.Second, cfg.IP.Duration)
	assert.Equal(t, 5*time.Minute, cfg.IP.BlockDuration)
	assert.Equal(t, 10, cfg.Token.Requests)
	assert.Equal(t, "8080", cfg.ServerPort)
}

func TestRedisConfig_Address(t *testing.T) {
	cfg := RedisConfig{Host: "localhost", Port: "6379"}
	assert.Equal(t, "localhost:6379", cfg.Address())
}
