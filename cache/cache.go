package cache

import (
	"time"
)

type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, data []byte) error
	SetWithTTL(key string, data []byte, ttl time.Duration) error
}
