package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/zhenianik/grpcApiTest/pkg/api"
	"time"
)

type RedisCache struct {
	host    string
	db      int
	expires time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) *RedisCache {
	return &RedisCache{
		host:    host,
		db:      db,
		expires: exp,
	}
}

func (cache *RedisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *RedisCache) Set(ctx context.Context, key string, user *api.UserList) error {
	client := cache.getClient()

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	client.Set(ctx, key, json, cache.expires*time.Second)

	return nil
}

func (cache *RedisCache) Get(ctx context.Context, key string) (*api.UserList, error) {
	client := cache.getClient()

	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return &api.UserList{}, err
	}

	users := api.UserList{}
	err = json.Unmarshal([]byte(val), &users)
	if err != nil {
		return &api.UserList{}, err
	}

	return &users, nil
}
