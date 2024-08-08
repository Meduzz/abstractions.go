package caching_test

import (
	"context"
	"testing"
	"time"

	"github.com/Meduzz/abstractions.go/internal/local/caching"
)

type (
	testData struct {
		value string
	}
)

func TestLocalCache(t *testing.T) {
	subject := caching.NewCache[testData](10 * time.Millisecond)
	data1 := &testData{"data1"}
	data2 := &testData{"data2"}
	ctx := context.Background()

	t.Run("happy camper", func(t *testing.T) {
		err := subject.Write(ctx, "data1", data1)

		if err != nil {
			t.Errorf("error was not nil, %v", err)
		}

		result, err := subject.Read(ctx, "data1")

		if err != nil {
			t.Errorf("error was not nil, %v", err)
		}

		if result.value != data1.value {
			t.Errorf("result was not data1 but %s", result.value)
		}
	})

	t.Run("expires stuff", func(t *testing.T) {
		err := subject.Write(ctx, "data2", data2)

		if err != nil {
			t.Errorf("error was not nil, %v", err)
		}

		timer := time.After(15 * time.Millisecond)

		<-timer

		result, err := subject.Read(ctx, "data2")

		if err == nil {
			t.Error("error was nil")
		}

		if result != nil {
			t.Error("result was not nil")
		}
	})
}
