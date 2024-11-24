package caching

import (
	"context"
	"errors"
	"time"

	"github.com/Meduzz/abstractions.go/lib"
	"github.com/Meduzz/abstractions.go/lib/specific"
	"github.com/go-redis/redis/v8"
)

type (
	abstraction struct {
		config   *specific.RedisConfig
		ttl      time.Duration
		eviction lib.Eviction
		name     string
	}
)

func NewCaching(config *specific.RedisConfig, eviction lib.Eviction, ttl time.Duration, name string) lib.CacheStorageDelegate {
	return &abstraction{
		config:   config,
		ttl:      ttl,
		name:     name,
		eviction: eviction,
	}
}

func (a *abstraction) Write(ctx context.Context, key string, data []byte) error {
	result := a.config.Redis().SetEX(ctx, a.config.Prefix(a.name, key), data, a.ttl)

	return result.Err()
}

func (a *abstraction) Read(ctx context.Context, key string) ([]byte, error) {
	result := a.config.Redis().GetEx(ctx, a.config.Prefix(a.name, key), a.ttl)

	bs, err := result.Bytes()

	if err != nil {
		if errors.Is(redis.Nil, err) {
			return nil, lib.ErrKeyNotFound
		}

		return nil, err
	}

	if a.eviction == lib.EvictionRead {
		err := a.config.Redis().Expire(ctx, a.config.Prefix(a.name, key), a.ttl).Err()

		if err != nil {
			return nil, err
		}
	}

	return bs, nil
}

func (a *abstraction) Delete(ctx context.Context, key string) error {
	return a.config.Redis().Del(ctx, a.config.Prefix(a.name, key)).Err()
}
