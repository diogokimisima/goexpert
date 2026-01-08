package storage

import (
	"context"
	"sync"
	"time"
)

type MemoryStorage struct {
	data      map[string]*entry
	mu        sync.RWMutex
	stopClean chan struct{}
}

type entry struct {
	value      int64
	expiration time.Time
}

func NewMemoryStorage() *MemoryStorage {
	m := &MemoryStorage{
		data:      make(map[string]*entry),
		stopClean: make(chan struct{}),
	}

	go m.cleanupExpired()

	return m
}

func (m *MemoryStorage) Increment(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	e, exists := m.data[key]

	if !exists || e.expiration.Before(now) {
		m.data[key] = &entry{
			value:      1,
			expiration: now.Add(expiration),
		}
		return 1, nil
	}

	e.value++
	return e.value, nil
}

func (m *MemoryStorage) Get(ctx context.Context, key string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	e, exists := m.data[key]
	if !exists || e.expiration.Before(time.Now()) {
		return 0, nil
	}

	return e.value, nil
}

func (m *MemoryStorage) SetBlock(ctx context.Context, key string, duration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	blockKey := key + ":blocked"
	m.data[blockKey] = &entry{
		value:      1,
		expiration: time.Now().Add(duration),
	}

	return nil
}

func (m *MemoryStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	blockKey := key + ":blocked"
	e, exists := m.data[blockKey]
	if !exists || e.expiration.Before(time.Now()) {
		return false, nil
	}

	return e.value == 1, nil
}

func (m *MemoryStorage) Close() error {
	close(m.stopClean)
	return nil
}

func (m *MemoryStorage) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mu.Lock()
			now := time.Now()
			for key, e := range m.data {
				if e.expiration.Before(now) {
					delete(m.data, key)
				}
			}
			m.mu.Unlock()
		case <-m.stopClean:
			return
		}
	}
}
