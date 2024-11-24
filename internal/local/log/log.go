package log

import (
	"context"

	"github.com/Meduzz/abstractions.go/lib"
	"github.com/Meduzz/helper/fp/slice"
)

type (
	localLog struct {
		storage []*lib.WorkItem
	}
)

func NewLocalLog() lib.WorkLogDelegate {
	return &localLog{make([]*lib.WorkItem, 0)}
}

func (l *localLog) Append(ctx context.Context, work *lib.WorkItem) error {
	l.storage = append(l.storage, work)

	return nil
}

func (l *localLog) Size(ctx context.Context) (int64, error) {
	return int64(len(l.storage)), nil
}

func (l *localLog) Fetch(ctx context.Context) (*lib.WorkItem, error) {
	work := slice.Head(l.storage)
	l.storage = slice.Tail(l.storage)

	return work, nil
}
