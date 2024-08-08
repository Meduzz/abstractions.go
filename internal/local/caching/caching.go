package caching

import (
	"context"
	"time"

	"github.com/Meduzz/abstractions.go/lib"
)

type (
	cacheValue[T any] struct {
		expire time.Time
		data   *T
	}
	localCache[T any] struct {
		ttl     time.Duration
		storage map[string]*cacheValue[T]
	}
)

func NewCache[T any](ttl time.Duration) lib.CacheAbstraction[T] {
	storage := make(map[string]*cacheValue[T])
	return &localCache[T]{ttl, storage}
}

func (l *localCache[T]) Write(ctx context.Context, key string, data *T) error {
	expires := time.Now().Add(l.ttl)
	item := &cacheValue[T]{
		expire: expires,
		data:   data,
	}

	l.storage[key] = item

	return nil
}

func (l *localCache[T]) Read(ctx context.Context, key string) (*T, error) {
	item, ok := l.storage[key]

	if !ok {
		return nil, lib.ErrKeyNotFound
	}

	if item.expire.Before(time.Now()) {
		delete(l.storage, key)
		return nil, lib.ErrKeyNotFound
	}

	item.expire = time.Now().Add(l.ttl)

	return item.data, nil
}
