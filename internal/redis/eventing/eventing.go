package eventing

import (
	"context"
	"encoding/json"

	"github.com/Meduzz/abstractions.go/lib"
	"github.com/go-redis/redis/v8"
)

type (
	eventingAbstraction struct {
		config *lib.RedisConfig
	}
)

func NewEventing(config *lib.RedisConfig) lib.EventingAbstraction {
	return &eventingAbstraction{config}
}

func (e *eventingAbstraction) Publish(ctx context.Context, topic string, payload interface{}) error {
	bs, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	return e.config.Redis().Publish(ctx, topic, bs).Err()
}

func (e *eventingAbstraction) Subscribe(ctx context.Context, pattern string, handler func(lib.Context)) {
	sub := e.config.Redis().PSubscribe(ctx, pattern)
	events := sub.Channel()

	go startEventHandler(events, sub, e.config, pattern, handler)
}

func startEventHandler(topic <-chan *redis.Message, sub *redis.PubSub, cfg *lib.RedisConfig, pattern string, handler func(lib.Context)) {
	for event := range topic {
		context := &eventingContext{
			sub:     sub,
			msg:     event,
			cfg:     cfg,
			pattern: pattern,
		}

		handler(context)
	}
}
