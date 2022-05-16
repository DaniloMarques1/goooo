package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache[T any] struct {
	redisConn *redis.Client
}

func NewRedisCache[T any](redisConn *redis.Client) *RedisCache[T] {
	return &RedisCache[T]{redisConn: redisConn}
}

func (rc *RedisCache[T]) SaveToCache(key string, value T, duration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = rc.redisConn.SetEX(
		context.Background(),
		key,
		string(bytes),
		duration,
	).Err()

	if err != nil {
		log.Printf("Error = %v\n", err)
		return err
	}
	return nil
}

// TODO maybe i should return T instead of receiving returnedValue as parameter
func (rs *RedisCache[T]) GetFromCache(key string, returnedValue T) error {
	str, err := rs.redisConn.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(str), returnedValue); err != nil {
		return err
	}
	return nil
}
