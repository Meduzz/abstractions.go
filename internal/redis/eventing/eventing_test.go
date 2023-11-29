package eventing_test

import (
	"context"
	"testing"

	"github.com/Meduzz/abstractions.go/internal/redis/eventing"
	"github.com/Meduzz/abstractions.go/lib"
	"github.com/Meduzz/helper/rudis"
)

type (
	testdata struct {
		Message string `json:"message"`
	}
)

func TestEventing(t *testing.T) {
	conn := rudis.Connect()
	cfg := lib.NewRedisConfig(conn, "")
	ctx := context.Background()
	result := make(chan string, 10)
	subject := eventing.NewEventing(cfg)
	subject.Subscribe(ctx, "test.*", listener(result))

	defer conn.Close()

	t.Run("test 1", func(t *testing.T) {
		subject.Publish(ctx, "test.1", &testdata{"test 1"})

		first := <-result

		if first != "test 1" {
			t.Error("first was not test 1", first)
		}
	})
}

func listener(gossip chan string) func(lib.Context) {
	return func(ctx lib.Context) {
		event := &testdata{}
		err := ctx.JSON(event)

		if err != nil {
			panic(err)
		}

		gossip <- event.Message
	}
}
