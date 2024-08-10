package caching_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Meduzz/abstractions.go/codec"
	"github.com/Meduzz/abstractions.go/internal/redis/caching"
	"github.com/Meduzz/abstractions.go/lib"
	"github.com/Meduzz/abstractions.go/lib/specific"
	"github.com/Meduzz/helper/rudis"
)

type (
	testdata struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
)

func TestCaching(t *testing.T) {
	conn := rudis.Connect()
	cfg := specific.NewRedisConfig(conn, "")
	subject := caching.NewCaching(cfg, codec.NewJsonCodec[testdata](), time.Second, "cache")
	data := &testdata{"project", 1}
	ctx := context.Background()

	defer conn.Close()

	t.Run("test write, extend and read from cache, incl expired cache", func(t *testing.T) {
		err := subject.Write(ctx, "cache.test1", data)

		if err != nil {
			t.Error("writing to cache threw error", err)
		}

		second, err := subject.Read(ctx, "cache.test1")

		if err != nil {
			t.Error("reading from cache threw error", err)
		}

		if data.Name != second.Name && data.Age != second.Age {
			t.Error("data was not equal to the result", data, second)
		}

		<-time.After(time.Second)
		_, err = subject.Read(ctx, "cache.test1")

		if err != nil {
			if !errors.Is(lib.ErrKeyNotFound, err) {
				t.Error("ErrKeyNotFound return as expected", err)
			}
		} else {
			t.Error("no error was thrown", data)
		}
	})

	t.Run("delete a cache item", func(t *testing.T) {
		err := subject.Write(ctx, "cache.test1", data)

		if err != nil {
			t.Error("writing to cache threw error", err)
		}

		second, err := subject.Read(ctx, "cache.test1")

		if err != nil {
			t.Error("reading from cache threw error", err)
		}

		if data.Name != second.Name && data.Age != second.Age {
			t.Error("data was not equal to the result", data, second)
		}

		err = subject.Del(ctx, "cache.test1")

		if err != nil {
			t.Error("deleting key from cache threw error", err)
		}

		_, err = subject.Read(ctx, "cache.test1")

		if err != nil {
			if !errors.Is(lib.ErrKeyNotFound, err) {
				t.Error("ErrKeyNotFound return as expected", err)
			}
		} else {
			t.Error("no error was thrown", data)
		}
	})
}
