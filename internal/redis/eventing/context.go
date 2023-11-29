package eventing

import (
	"context"
	"encoding/json"

	"github.com/Meduzz/abstractions.go/lib"
	"github.com/go-redis/redis/v8"
)

type (
	eventingContext struct {
		sub     *redis.PubSub
		msg     *redis.Message
		cfg     *lib.RedisConfig
		pattern string
	}
)

func newContext(sub *redis.PubSub, msg *redis.Message, cfg *lib.RedisConfig, pattern string) lib.Context {
	return &eventingContext{
		sub:     sub,
		msg:     msg,
		cfg:     cfg,
		pattern: pattern,
	}
}

func (e *eventingContext) Topic() string {
	return e.msg.Channel
}

func (e *eventingContext) Bytes() []byte {
	return []byte(e.msg.Pattern)
}

func (e *eventingContext) JSON(into interface{}) error {
	return json.Unmarshal([]byte(e.msg.Payload), into)
}

func (e *eventingContext) Forward(ctx context.Context, topic string) error {
	return e.cfg.Redis().Publish(ctx, topic, e.msg.Payload).Err()
}

func (e *eventingContext) Unsubscribe(ctx context.Context) {
	e.sub.Unsubscribe(ctx, e.pattern)
}
