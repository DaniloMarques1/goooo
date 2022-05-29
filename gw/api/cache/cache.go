package cache

import "time"

type Cache interface {
	SaveToCache(key string, value []byte, duration time.Duration) error
	GetFromCache(key string) ([]byte, error)
}
