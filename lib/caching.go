package lib

import (
	"context"
	"time"
)

type (
	// CacheAbstraction abstrac away caching over a redis db.
	CacheAbstraction[T any] interface {
		// Write - write the provided data to the cache on the provided key.
		Write(ctx context.Context, key string, duration time.Duration, data *T) error
		// Read - read data from the provided key, if exists otherwise return error.
		Read(ctx context.Context, key string) (*T, error)
		// ReadAndExtend - read data from the provided key, extend if it exists, otherwise return error.
		ReadAndExtend(ctx context.Context, key string, extend time.Duration) (*T, error)
	}
)
