package log

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/Meduzz/abstractions.go/lib"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	testData struct {
		Message string `json:"message"`
	}
)

func TestDbms(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	subject, err := NewDbmsWorkLogDelegate(db)

	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	work1 := createTestData(&testData{"work1"})
	work2 := createTestData(&testData{"work2"})

	t.Run("happy camper", func(t *testing.T) {
		err = subject.Append(ctx, work1)

		if err != nil {
			t.Error(err)
		}

		err = subject.Append(ctx, work2)

		if err != nil {
			t.Error(err)
		}

		size, err := subject.Size(ctx)

		if err != nil {
			t.Error(err)
		}

		if size != 2 {
			t.Errorf("size was not 2 but %d", size)
		}

		result1, err := subject.Fetch(ctx)

		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(work1.Work, result1.Work) {
			t.Errorf("result.Work was not work1 but %s", result1.Work)
		}

		size, err = subject.Size(ctx)

		if err != nil {
			t.Error(err)
		}

		if size != 1 {
			t.Errorf("size was not 1 but %d", size)
		}

		result2, err := subject.Fetch(ctx)

		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(work2.Work, result2.Work) {
			t.Errorf("result.Work was not work2 but %s", result1.Work)
		}

		size, err = subject.Size(ctx)

		if err != nil {
			t.Error(err)
		}

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
