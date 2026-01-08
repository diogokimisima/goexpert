package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_IncrementAndGet(t *testing.T) {
	store := NewMemoryStorage()
	defer store.Close()
	ctx := context.Background()

	count, err := store.Increment(ctx, "test-key", time.Second)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	count, err = store.Increment(ctx, "test-key", time.Second)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	count, err = store.Get(ctx, "test-key")
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestMemoryStorage_BlockAndIsBlocked(t *testing.T) {
	store := NewMemoryStorage()
	defer store.Close()
	ctx := context.Background()

	blocked, err := store.IsBlocked(ctx, "test-key")
	require.NoError(t, err)
	assert.False(t, blocked)

	err = store.SetBlock(ctx, "test-key", time.Second)
	require.NoError(t, err)

	blocked, err = store.IsBlocked(ctx, "test-key")
	require.NoError(t, err)
	assert.True(t, blocked)
}
