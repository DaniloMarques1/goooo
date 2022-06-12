package cache

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (rc *RedisCache) SaveToCache(key string, value []byte, duration time.Duration) error {
	err := rc.client.SetEX(
		context.Background(),
		key,
		string(value),
		duration,
	).Err()

	if err != nil {
		log.Printf("Error = %v\n", err)
		return err
	}
	return nil
}

func (rs *RedisCache) GetFromCache(key string) ([]byte, error) {
	str, err := rs.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	return []byte(str), nil
}
