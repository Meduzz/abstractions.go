package log_test

import (
	"context"
	"testing"

	"github.com/Meduzz/abstractions.go/internal/local/log"
)

type (
	testData struct {
		message string
	}
)

func TestLocalLog(t *testing.T) {
	ctx := context.Background()
	subject := log.NewLocalLog[testData]()
	work1 := &testData{"work1"}
	work2 := &testData{"work2"}

	t.Run("happy camper", func(t *testing.T) {
		subject.Append(ctx, work1)
		subject.Append(ctx, work2)

		size, _ := subject.Size(ctx)

		if size != 2 {
			t.Errorf("size was not 2 but %d", size)
		}

		result1, _ := subject.Fetch(ctx)

		if result1.message != work1.message {
			t.Errorf("result.message was not work1 but %s", result1.message)
		}

		size, _ = subject.Size(ctx)

		if size != 1 {
			t.Errorf("size was not 1 but %d", size)
		}

		result2, _ := subject.Fetch(ctx)

		if result2.message != work2.message {
			t.Errorf("result.message was not work2 but %s", result1.message)
		}

		size, _ = subject.Size(ctx)

		if size != 0 {
			t.Errorf("size was not 0 but %d", size)
		}
	})
}
