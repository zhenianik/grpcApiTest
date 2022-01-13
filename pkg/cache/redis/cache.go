package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client  *redis.Client
	host    string
	db      int
	expires time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       db,
	})

	return &RedisCache{
		client:  client,
		host:    host,
		db:      db,
		expires: exp,
	}
}

func (c *RedisCache) Set(ctx context.Context, key string, user []byte) error {
	c.client.Set(ctx, key, user, c.expires*time.Second)
	return nil
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}
