package caching

import (
	"context"
	"time"

	"github.com/Meduzz/abstractions.go/internal/interval"
	"github.com/Meduzz/abstractions.go/lib"
)

type (
	cacheValue struct {
		expire time.Time
		data   []byte
	}
	localCache struct {
		eviction lib.Eviction
		ttl      time.Duration
		storage  map[string]*cacheValue
	}
)

func NewCache(eviction lib.Eviction, ttl time.Duration) lib.CacheStorageDelegate {
	storage := make(map[string]*cacheValue)

	interval.OnInterval(5*time.Second, func() {
		for k, v := range storage {
			if v.expire.Before(time.Now()) {
				delete(storage, k)
			}
		}
	})

	return &localCache{eviction, ttl, storage}
}

func (l *localCache) Write(ctx context.Context, key string, data []byte) error {
	expires := time.Now().Add(l.ttl)
	item := &cacheValue{
		expire: expires,
		data:   data,
	}

	l.storage[key] = item

	return nil
}

func (l *localCache) Read(ctx context.Context, key string) ([]byte, error) {
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

func (l *localCache) Delete(ctx context.Context, key string) error {
	delete(l.storage, key)

	return nil
}
