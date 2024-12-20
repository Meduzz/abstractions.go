package log_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/Meduzz/abstractions.go/internal/local/log"
	"github.com/Meduzz/abstractions.go/lib"
)

type (
	testData struct {
		Message string `json:"message"`
	}
)

func TestLocalLog(t *testing.T) {
	ctx := context.Background()
	subject := log.NewLocalLog()
	work1 := createTestData(&testData{"work1"})
	work2 := createTestData(&testData{"work2"})

	t.Run("happy camper", func(t *testing.T) {
		subject.Append(ctx, work1)
		subject.Append(ctx, work2)

		size, _ := subject.Size(ctx)

		if size != 2 {
			t.Errorf("size was not 2 but %d", size)
		}

		result1, _ := subject.Fetch(ctx)

		if !bytes.Equal(work1.Work, result1.Work) {
			t.Errorf("result.Work was not work1 but %s", result1.Work)
		}

		size, _ = subject.Size(ctx)

		if size != 1 {
			t.Errorf("size was not 1 but %d", size)
		}

		result2, _ := subject.Fetch(ctx)

		if !bytes.Equal(work2.Work, result2.Work) {
			t.Errorf("result.Work was not work2 but %s", result1.Work)
		}

		size, _ = subject.Size(ctx)

		if size != 0 {
			t.Errorf("size was not 0 but %d", size)
		}
	})
}

func createTestData(data *testData) *lib.WorkItem {
	work := &lib.WorkItem{}

	bs, _ := json.Marshal(data)

	work.Kind = "work"
	work.Work = bs

	return work
}
