package caching_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Meduzz/abstractions.go/codec"
	root "github.com/Meduzz/abstractions.go/internal/redis"
	"github.com/Meduzz/abstractions.go/internal/redis/caching"
	"github.com/Meduzz/abstractions.go/lib"
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
	cfg := root.NewRedisConfig(conn, "")

	defer conn.Close()

	t.Run("test write, extend and read from cache, incl expired cache", func(t *testing.T) {
		module := caching.NewCaching(cfg, codec.NewJsonCodec[testdata](), time.Second)
		data := &testdata{"project", 1}
		ctx := context.Background()

		err := module.Write(ctx, "cache.test1", data)

		if err != nil {
			t.Error("writing to cache threw error", err)
		}

		second, err := module.Read(ctx, "cache.test1")

		if err != nil {
			t.Error("reading from cache threw error", err)
		}

		if data.Name != second.Name && data.Age != second.Age {
			t.Error("data was not equal to the result", data, second)
		}

		<-time.After(time.Second)
		_, err = module.Read(ctx, "cache.test1")

		if err != nil {
			if !errors.Is(lib.ErrKeyNotFound, err) {
				t.Error("ErrKeyNotFound return as expected", err)
			}
		} else {
			t.Error("no error was thrown", data)
		}
	})
}
