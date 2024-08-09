package log

import (
	"context"

	"github.com/Meduzz/abstractions.go/lib"
	"github.com/Meduzz/helper/fp/slice"
)

type (
	localLog[T any] struct {
		storage []*T
	}
)

func NewLocalLog[T any]() lib.LogAbstraction[T] {
	return &localLog[T]{make([]*T, 0)}
}

func (l *localLog[T]) Append(ctx context.Context, work *T) error {
	l.storage = append(l.storage, work)

	return nil
}

func (l *localLog[T]) Size(ctx context.Context) (int64, error) {
	return int64(len(l.storage)), nil
}

func (l *localLog[T]) Fetch(ctx context.Context) (*T, error) {
	work := slice.Head(l.storage)
	l.storage = slice.Tail(l.storage)

	return work, nil
}
