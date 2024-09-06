package cache

import (
	"sync"
	"time"
)

type Item struct {
	Expiry time.Time
	Value  []byte
}

const (
	NeverExpires time.Duration = -1
	DefaultTTL   time.Duration = 0
)

func (i *Item) Expired() bool {
	return i.Expiry.Before(time.Now())
}

type MemCache struct {
	defaultTTL time.Duration
	items      map[string]Item
	mu         sync.RWMutex
}

func NewMemCache(defaultTTL time.Duration) *MemCache {
	return &MemCache{
		defaultTTL: defaultTTL,
		items:      make(map[string]Item),
	}
}

// Set sets the value for the key with the default TTL
func (m *MemCache) Set(key string, value []byte) error {
	return m.SetWithTTL(key, value, m.defaultTTL)
}

// SetWithTTL sets the value for the key with the specified defaultTTL
func (m *MemCache) SetWithTTL(key string, value []byte, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[key] = Item{
		Expiry: time.Now().Add(ttl),
		Value:  value,
	}
	return nil
}

// Get returns the value for the key if it exists and is not Expired
func (m *MemCache) Get(key string) ([]byte, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, ok := m.items[key]
	if !ok || item.Expired() {
		return nil, false
	}

	return item.Value, true
}
