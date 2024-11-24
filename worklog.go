package abstractions

import (
	"context"
	"time"

	"github.com/Meduzz/abstractions.go/lib"
)

type (
	WorkLogHandler interface {
		Work(*lib.WorkItem)
		Error(error)
	}

	WorkLog struct {
		delegate lib.WorkLogDelegate
		handler  WorkLogHandler
		interval time.Duration
		ticker   *time.Ticker
	}
)

func NewWorkLog(interval time.Duration, delegate lib.WorkLogDelegate, handler WorkLogHandler) *WorkLog {
	return &WorkLog{
		delegate: delegate,
		handler:  handler,
		interval: interval,
	}
}

func (w *WorkLog) Start() {
	w.ticker = time.NewTicker(w.interval)

	go func() {
		ctx := context.Background()
		for range w.ticker.C {
			size, _ := w.delegate.Size(ctx)

			if size > 0 {
				item, err := w.delegate.Fetch(ctx)

				if err != nil {
					w.handler.Error(err)
				} else {
					w.handler.Work(item)
				}
			}
		}
	}()
}

func (w *WorkLog) Stop() {
	w.ticker.Stop()
}

func (w *WorkLog) Append(work *lib.WorkItem) error {
	ctx := context.Background()

	return w.delegate.Append(ctx, work)
}
