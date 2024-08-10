package eventing

import (
	"context"

	"github.com/Meduzz/abstractions.go/lib"
	"github.com/Meduzz/abstractions.go/lib/vendor"
	"github.com/go-redis/redis/v8"
)

type (
	eventingAbstraction[T any] struct {
		config *vendor.RedisConfig
		topic  string
		codec  lib.Codec[T]
	}
)

func NewEventing[T any](topic string, codec lib.Codec[T], config *vendor.RedisConfig) lib.EventingAbstraction[T] {
	return &eventingAbstraction[T]{config, topic, codec}
}

func (e *eventingAbstraction[T]) Publish(ctx context.Context, payload *T) error {
	bs, err := e.codec.Encode(payload)

	if err != nil {
		return err
	}

	return e.config.Redis().Publish(ctx, e.topic, bs).Err()
}

func (e *eventingAbstraction[T]) Subscribe(ctx context.Context, handler func(*T)) {
	sub := e.config.Redis().PSubscribe(ctx, e.topic)
	events := sub.Channel()

	go e.startEventHandler(events, handler)
}

func (e *eventingAbstraction[T]) startEventHandler(topic <-chan *redis.Message, handler func(*T)) {
	for event := range topic {
		payload, err := e.codec.Decode([]byte(event.Payload))

		if err != nil {
			// TODO spew this onto stdout somehow?
			continue
		}

		handler(payload)
	}
}
