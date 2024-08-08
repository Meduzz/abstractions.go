package eventing_test

import (
	"context"
	"testing"

	"github.com/Meduzz/abstractions.go/internal/local/eventing"
)

type (
	testData struct {
		message string
	}
)

func TestLocalEventing(t *testing.T) {
	transport := make(chan *testData)
	subject := eventing.NewLocalEventing(transport)
	results := make(chan string)
	ctx := context.Background()

	// setup a listener
	subject.Subscribe(ctx, func(td *testData) {
		results <- td.message
	})

	t.Run("happy campers", func(t *testing.T) {
		data := &testData{"test1"}

		err := subject.Publish(ctx, data)

		if err != nil {
			t.Errorf("error was not nil, was %v", err)
		}

		result := <-results

		if result != "test1" {
			t.Errorf("result was not test1 but %s", result)
		}
	})

	close(results)
	close(transport)
}
