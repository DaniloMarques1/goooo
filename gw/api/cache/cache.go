package cache

import "time"

type Cache[T any] interface {
	SaveToCache(key string, value T, duration time.Duration) error
	GetFromCache(key string, returnedValue T) error
}
