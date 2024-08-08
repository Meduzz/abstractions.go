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
		ttl    time.Duration
		name   string
	}
)

func NewCSRFAbstraction(config *root.RedisConfig, ttl time.Duration, name string) lib.CSRFAbstraction {
	return &abstraction{
		config: config,
		ttl:    ttl,
		name:   name,
	}
}

func (a *abstraction) Generate(ctx context.Context) (*lib.CSRFToken, error) {
	key := hashing.Token()
	value := hashing.Secret()

	token := &lib.CSRFToken{
		Key:   key,
		Value: value,
	}

	err := a.config.Redis().Set(ctx, a.config.Prefix(a.name, key), value, a.ttl).Err()

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (a *abstraction) Verify(ctx context.Context, key, value string) (bool, error) {
	data, err := a.config.Redis().GetDel(ctx, a.config.Prefix(a.name, key)).Result()

	if err != nil {
		if errors.Is(redis.Nil, err) {
			return false, nil
		}

		return false, err
	}

	return data == value, nil
}
