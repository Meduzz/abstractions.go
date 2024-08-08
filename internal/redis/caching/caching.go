package caching

import (
	"context"
	"errors"
	"time"

	root "github.com/Meduzz/abstractions.go/internal/redis"
	"github.com/Meduzz/abstractions.go/lib"
	"github.com/go-redis/redis/v8"
)

type (
	abstraction[T any] struct {
		config *root.RedisConfig
		codec  lib.Codec[T]
		ttl    time.Duration
	}
)

func NewCaching[T any](config *root.RedisConfig, codec lib.Codec[T], ttl time.Duration) lib.CacheAbstraction[T] {
	return &abstraction[T]{
		config: config,
		codec:  codec,
		ttl:    ttl,
	}
}

func (a *abstraction[T]) Write(ctx context.Context, key string, data *T) error {
	bs, err := a.codec.Encode(data)

	if err != nil {
		return err
	}

	result := a.config.Redis().Set(ctx, a.config.Prefix(key), bs, a.ttl)

	return result.Err()
}

func (a *abstraction[T]) Read(ctx context.Context, key string) (*T, error) {
	result := a.config.Redis().Get(ctx, a.config.Prefix(key))

	bs, err := result.Bytes()

	if err != nil {
		if errors.Is(redis.Nil, err) {
			return nil, lib.ErrKeyNotFound
		}

		return nil, err
	}

	err = a.config.Redis().Set(ctx, a.config.Prefix(key), bs, a.ttl).Err()

	if err != nil {
		return nil, err
	}

	return a.codec.Decode(bs)
}
