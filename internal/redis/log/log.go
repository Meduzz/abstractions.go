package log

import (
	"context"
	"errors"

	root "github.com/Meduzz/abstractions.go/internal/redis"
	"github.com/Meduzz/abstractions.go/lib"
	"github.com/go-redis/redis/v8"
)

type (
	logAbstraction[T any] struct {
		config *root.RedisConfig
		codec  lib.Codec[T]
		name   string
	}
)

func NewRedisLog[T any](config *root.RedisConfig, codec lib.Codec[T], name string) lib.LogAbstraction[T] {
	fullName := config.Prefix(name)

	return &logAbstraction[T]{
		config: config,
		codec:  codec,
		name:   fullName,
	}
}

func (l *logAbstraction[T]) Append(ctx context.Context, work *T) error {
	bs, err := l.codec.Encode(work)

	if err != nil {
		return err
	}

	return l.config.Redis().RPush(ctx, l.name, bs).Err()
}

func (l *logAbstraction[T]) Size(ctx context.Context) (int64, error) {
	size, err := l.config.Redis().LLen(ctx, l.name).Result()

	if err != nil {
		if errors.Is(redis.Nil, err) {
			return 0, lib.ErrKeyNotFound
		}

		return 0, err
	}

	return size, nil
}

func (l *logAbstraction[T]) Fetch(ctx context.Context) (*T, error) {
	workBytes, err := l.config.Redis().LPop(ctx, l.name).Result()

	if err != nil {
		if errors.Is(redis.Nil, err) {
			return nil, lib.ErrKeyNotFound
		}

		return nil, err
	}

	return l.codec.Decode([]byte(workBytes))
}
