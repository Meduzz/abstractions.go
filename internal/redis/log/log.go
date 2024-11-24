package log

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Meduzz/abstractions.go/lib"
	"github.com/Meduzz/abstractions.go/lib/specific"
	"github.com/go-redis/redis/v8"
)

type (
	redisWorkLogDelegate struct {
		config *specific.RedisConfig
		name   string
	}
)

func NewRedisWorkLog(config *specific.RedisConfig, name string) lib.WorkLogDelegate {
	fullName := config.Prefix(name)

	return &redisWorkLogDelegate{
		config: config,
		name:   fullName,
	}
}

func (l *redisWorkLogDelegate) Append(ctx context.Context, work *lib.WorkItem) error {
	bs, err := json.Marshal(work)

	if err != nil {
		return err
	}

	return l.config.Redis().RPush(ctx, l.name, bs).Err()
}

func (l *redisWorkLogDelegate) Size(ctx context.Context) (int64, error) {
	size, err := l.config.Redis().LLen(ctx, l.name).Result()

	if err != nil {
		if errors.Is(redis.Nil, err) {
			return 0, lib.ErrKeyNotFound
		}

		return 0, err
	}

	return size, nil
}

func (l *redisWorkLogDelegate) Fetch(ctx context.Context) (*lib.WorkItem, error) {
	workBytes, err := l.config.Redis().LPop(ctx, l.name).Result()

	if err != nil {
		if errors.Is(redis.Nil, err) {
			return nil, lib.ErrKeyNotFound
		}

		return nil, err
	}

	work := &lib.WorkItem{}
	err = json.Unmarshal([]byte(workBytes), work)

	if err != nil {
		return nil, err
	}

	return work, nil
}
