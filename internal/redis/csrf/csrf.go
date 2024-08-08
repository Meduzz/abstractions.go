package csrf

import (
	"context"
	"errors"
	"time"

	root "github.com/Meduzz/abstractions.go/internal/redis"
	"github.com/Meduzz/abstractions.go/lib"
	"github.com/Meduzz/helper/hashing"
	"github.com/go-redis/redis/v8"
)

type (
	abstraction struct {
		config *root.RedisConfig
	}
)

func NewCSRFAbstraction(config *root.RedisConfig) lib.CSRFAbstraction {
	return &abstraction{
		config: config,
	}
}

func (a *abstraction) Generate(ctx context.Context, duration time.Duration) (*lib.CSRFToken, error) {
	key := hashing.Token()
	value := hashing.Secret()

	token := &lib.CSRFToken{
		Key:   key,
		Value: value,
	}

	err := a.config.Redis().Set(ctx, a.config.Prefix(key), value, duration).Err()

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (a *abstraction) Verify(ctx context.Context, key, value string) (bool, error) {
	data, err := a.config.Redis().Get(ctx, a.config.Prefix(key)).Result()

	if err != nil {
		if errors.Is(redis.Nil, err) {
			return false, nil
		}

		return false, err
	}

	err = a.config.Redis().Del(ctx, a.config.Prefix(key)).Err()

	return data == value, err
}
