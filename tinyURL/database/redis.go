package database

import (
	"time"
	"tinyURL/config"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(cfg *config.Config) (*RedisCache, error) {

	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

func (c *RedisCache) Set(key string, value string, expiration time.Duration) error {
	return c.client.Set(c.client.Context(), key, value, expiration).Err()
}

func (c *RedisCache) Get(key string) (string, error) {
	return c.client.Get(c.client.Context(), key).Result()
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}
