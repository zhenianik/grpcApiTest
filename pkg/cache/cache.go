package cache

import (
	"context"
)

//go:generate mockgen -source=cache.go -destination=./mock/cache.go -package=mock

type (
	Cache interface {
		Set(ctx context.Context, key string, user []byte) error
		Get(ctx context.Context, key string) ([]byte, error)
	}
)
