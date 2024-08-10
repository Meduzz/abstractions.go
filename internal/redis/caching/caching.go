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
	abstraction[T any] struct {
		config *specific.RedisConfig
		codec  lib.Codec[T]
		ttl    time.Duration
		name   string
	}
)

func NewCaching[T any](config *specific.RedisConfig, codec lib.Codec[T], ttl time.Duration, name string) lib.CacheAbstraction[T] {
	return &abstraction[T]{
		config: config,
		codec:  codec,
		ttl:    ttl,
		name:   name,
	}
}

func (a *abstraction[T]) Write(ctx context.Context, key string, data *T) error {
	bs, err := a.codec.Encode(data)

	if err != nil {
		return err
	}

	result := a.config.Redis().SetEX(ctx, a.config.Prefix(a.name, key), bs, a.ttl)

	return result.Err()
}

func (a *abstraction[T]) Read(ctx context.Context, key string) (*T, error) {
	result := a.config.Redis().GetEx(ctx, a.config.Prefix(a.name, key), a.ttl)

	bs, err := result.Bytes()

	if err != nil {
		if errors.Is(redis.Nil, err) {
			return nil, lib.ErrKeyNotFound
		}

		return nil, err
	}

	return a.codec.Decode(bs)
}

func (a *abstraction[T]) Del(ctx context.Context, key string) error {
	return a.config.Redis().Del(ctx, a.config.Prefix(a.name, key)).Err()
}
