package lib

import (
	"context"
)

type (
	// EventAbstraction - Send and receive events of type T.
	EventingAbstraction[T any] interface {
		// Publish - send a message to the specified topic.
		Publish(ctx context.Context, data *T) error
		// Subscribe - subscribe the handler to the provided pattern (psubscribe).
		Subscribe(ctx context.Context, handler func(*T))
	}
)
