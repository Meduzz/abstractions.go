package lib

import (
	"context"
)

type (
	Eviction string

	// CacheAbstraction abstract away caching.
	CacheStorageDelegate interface {
		// Write - write the provided data to the cache on the provided key.
		Write(ctx context.Context, key string, data []byte) error
		// Read - read data from the provided key (and extend it), if exists otherwise return error.
		Read(ctx context.Context, key string) ([]byte, error)
		// Delete - deletes a cache key.
		Delete(ctx context.Context, key string) error
	}
)

var (
	EvictionRead  = Eviction("read")
	EvictionWrite = Eviction("write")
)
