package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zhenianik/grpcApiTest/internal/controller/api"
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

func (cache *RedisCache) Set(ctx context.Context, key string, user *api.UserList) error {

	js, err := json.Marshal(user)
	if err != nil {
		return err
	}

	cache.client.Set(ctx, key, js, cache.expires*time.Second)

	return nil
}

func (cache *RedisCache) Get(ctx context.Context, key string) (*api.UserList, error) {

	val, err := cache.client.Get(ctx, key).Result()
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
