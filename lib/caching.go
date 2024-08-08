package lib

import (
	"context"
)

type (
	// CacheAbstraction abstract away caching.
	CacheAbstraction[T any] interface {
		// Write - write the provided data to the cache on the provided key.
		Write(ctx context.Context, key string, data *T) error
		// Read - read data from the provided key (and extend it), if exists otherwise return error.
		Read(ctx context.Context, key string) (*T, error)
	}
)
