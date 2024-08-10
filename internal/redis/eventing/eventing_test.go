package eventing_test

import (
	"context"
	"testing"

	"github.com/Meduzz/abstractions.go/codec"
	"github.com/Meduzz/abstractions.go/internal/redis/eventing"
	"github.com/Meduzz/abstractions.go/lib/vendor"
	"github.com/Meduzz/helper/rudis"
)

type (
	testdata struct {
		Message string `json:"message"`
	}
)

func TestEventing(t *testing.T) {
	conn := rudis.Connect()
	cfg := vendor.NewRedisConfig(conn, "")
	codec := codec.NewJsonCodec[testdata]()
	ctx := context.Background()
	result := make(chan string, 10)
	subject := eventing.NewEventing("test.1", codec, cfg)
	subject.Subscribe(ctx, listener(result))

	defer conn.Close()

	t.Run("test 1", func(t *testing.T) {
		subject.Publish(ctx, &testdata{"test 1"})

		first := <-result

		if first != "test 1" {
			t.Error("first was not test 1", first)
		}
	})
}

func listener(gossip chan string) func(*testdata) {
	return func(it *testdata) {
		gossip <- it.Message
	}
}
