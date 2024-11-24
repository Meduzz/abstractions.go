package lib

import (
	"context"
	"encoding/json"
)

type (
	WorkItem struct {
		Kind string          `json:"kind"`
		Work json.RawMessage `json:"work"`
	}

	WorkLogDelegate interface {
		// Append - append work to the queue
		Append(context.Context, *WorkItem) error
		// Size - fetch the size of the queue
		Size(context.Context) (int64, error)
		// Fetch - fetch the first item in the log and remove it
		Fetch(context.Context) (*WorkItem, error)
	}
)
