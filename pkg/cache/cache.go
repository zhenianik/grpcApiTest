package cache

import (
	"context"
)

//go:generate mockgen -source=cache.go -destination=./mock/cache.go -package=mock

type (
	Cache interface {
		Set(ctx context.Context, key string, user []byte) error
		Delete(ctx context.Context, key string) error
		Get(ctx context.Context, key string) ([]byte, error)
	}
)
