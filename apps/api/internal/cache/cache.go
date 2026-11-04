package cache

import (
	"context"
	"errors"
)

var ErrMiss = errors.New("cache miss")

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Put(ctx context.Context, key string, audio []byte) error
}
