package eventing

import (
	"context"

	"github.com/Meduzz/abstractions.go/lib"
)

type (
	localEventing[T any] struct {
		transport chan *T
	}
)

func NewLocalEventing[T any](transport chan *T) lib.EventingAbstraction[T] {
	return &localEventing[T]{transport}
}

func (l *localEventing[T]) Publish(ctx context.Context, data *T) error {
	l.transport <- data

	return nil
}

func (l *localEventing[T]) Subscribe(ctx context.Context, handler func(*T)) {
	go func() {
		for payload := range l.transport {
			handler(payload)
		}
	}()
}
