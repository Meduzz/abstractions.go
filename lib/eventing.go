package lib

import (
	"context"
)

type (
	// EventAbstraction - in order to not make assumptions, this became pretty low level.
	EventingAbstraction interface {
		// Publish - send a message to the specified topic.
		Publish(ctx context.Context, topic string, data interface{}) error
		// Subscribe - subscribe the handler to the provided pattern (psubscribe).
		Subscribe(ctx context.Context, pattern string, handler func(Context))
	}

	Context interface {
		Topic() string
		Bytes() []byte
		JSON(interface{}) error
		Forward(context.Context, string) error
		Unsubscribe(context.Context)
	}
)
