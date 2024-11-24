package log_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/Meduzz/abstractions.go/internal/redis/log"
	"github.com/Meduzz/abstractions.go/lib"
	"github.com/Meduzz/abstractions.go/lib/specific"
	"github.com/Meduzz/helper/rudis"
)

type (
	testData struct {
		Message string `json:"message"`
	}
)

func TestRedisLog(t *testing.T) {
	conn := rudis.Connect()
	config := specific.NewRedisConfig(conn, "test1")
	subject := log.NewRedisWorkLog(config, "testData")
	ctx := context.Background()
	work1 := createTestData(&testData{"work1"})
	work2 := createTestData(&testData{"work2"})

	t.Run("happy camper", func(t *testing.T) {
		err := subject.Append(ctx, work1)

		if err != nil {
			t.Errorf("appending work threw error %v", err)
		}

		err = subject.Append(ctx, work2)

		if err != nil {
			t.Errorf("appending work threw error %v", err)
		}

		size, err := subject.Size(ctx)

		if err != nil {
			t.Errorf("calculating size threw error %v", err)
		}

		if size != 2 {
			t.Errorf("size was not 2 but %d", size)
		}

		result1, err := subject.Fetch(ctx)

		if err != nil {
			t.Errorf("fetching work threw error %v", err)
		}

		if !bytes.Equal(work1.Work, result1.Work) {
			t.Errorf("result.work was not work1 but %s", result1.Work)
		}

		size, err = subject.Size(ctx)

		if err != nil {
			t.Errorf("calculating size threw error %v", err)
		}

		if size != 1 {
			t.Errorf("size was not 1 but %d", size)
		}

		result2, err := subject.Fetch(ctx)

		if err != nil {
			t.Errorf("fetching work threw error %v", err)
		}

		if !bytes.Equal(work2.Work, result2.Work) {
			t.Errorf("result.work was not work2 but %s", result1.Work)
		}

		size, err = subject.Size(ctx)

		if err != nil {
			t.Errorf("calculating size threw error %v", err)
		}

		if size != 0 {
			t.Errorf("size was not 0 but %d", size)
		}

		result3, err := subject.Fetch(ctx)

		if err == nil {
			t.Error("fetching work on empty queue did not throw error")
		}

		if result3 != nil {
			t.Errorf("result was not nil but %s", result1.Work)
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
