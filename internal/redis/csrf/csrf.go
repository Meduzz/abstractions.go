package csrf

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
		config *specific.RedisConfig
		ttl    time.Duration
		name   string
	}
)

func NewRedisCSRFStorageDelegate(config *specific.RedisConfig, ttl time.Duration, name string) lib.CSRFStorageDelegate {
	return &abstraction{
		config: config,
		ttl:    ttl,
		name:   name,
	}
}

func (a *abstraction) Store(ctx context.Context, token *lib.CSRFToken) error {
	err := a.config.Redis().Set(ctx, a.config.Prefix(a.name, token.Key), token.Value, a.ttl).Err()

	if err != nil {
		return err
	}

	return nil
}

func (a *abstraction) Verify(ctx context.Context, token *lib.CSRFToken) (bool, error) {
	data, err := a.config.Redis().GetDel(ctx, a.config.Prefix(a.name, token.Key)).Result()

	if err != nil {
		if errors.Is(redis.Nil, err) {
			return false, nil
		}

		return false, err
	}

	return data == token.Value, nil
}
