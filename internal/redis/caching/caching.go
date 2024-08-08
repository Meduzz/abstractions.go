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
		name   string
	}
)

func NewCaching[T any](config *root.RedisConfig, codec lib.Codec[T], ttl time.Duration, name string) lib.CacheAbstraction[T] {
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
