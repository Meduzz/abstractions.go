package lib

import "context"

type (
	// LogAbstraction - can be though of as a work queue.
	LogAbstraction[T any] interface {
		// Append - append work to the queue
		Append(context.Context, *T) error
		// Size - fetch the size of the queue
		Size(context.Context) (int64, error)
		// Fetch - fetch the first item in the log and remove it
		Fetch(context.Context) (*T, error)
	}
)
