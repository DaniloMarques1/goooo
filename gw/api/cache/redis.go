package cache

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	redisConn *redis.Client
}

func NewRedisCache(redisConn *redis.Client) *RedisCache {
	return &RedisCache{redisConn: redisConn}
}

func (rc *RedisCache) SaveToCache(key string, value []byte, duration time.Duration) error {
	err := rc.redisConn.SetEX(
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
	str, err := rs.redisConn.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	return []byte(str), nil
}
