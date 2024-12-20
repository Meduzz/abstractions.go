package caching_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Meduzz/abstractions.go/internal/redis/caching"
	"github.com/Meduzz/abstractions.go/lib"
	"github.com/Meduzz/abstractions.go/lib/specific"
	"github.com/Meduzz/helper/rudis"
)

func TestCaching(t *testing.T) {
	conn := rudis.Connect()
	cfg := specific.NewRedisConfig(conn, "")
	writeSubject := caching.NewCaching(cfg, lib.EvictionWrite, time.Second, "cache")
	readSubject := caching.NewCaching(cfg, lib.EvictionRead, time.Second, "cache")
	data := []byte("Im data")
	ctx := context.Background()

	defer conn.Close()

	t.Run("test caching with writeEviction", func(t *testing.T) {
		err := writeSubject.Write(ctx, "cache.test1", data)

		if err != nil {
			t.Error("writing to cache threw error", err)
		}

		second, err := writeSubject.Read(ctx, "cache.test1")

		if err != nil {
			t.Error("reading from cache threw error", err)
		}

		if !bytes.Equal(data, second) {
			t.Error("data was not equal to the result", data, second)
		}

		<-time.After(time.Second)
		_, err = writeSubject.Read(ctx, "cache.test1")

		if err != nil {
			if !errors.Is(err, lib.ErrKeyNotFound) {
				t.Error("ErrKeyNotFound was not returned as expected", err)
			}
		} else {
			t.Error("no error was thrown", data)
		}
	})

	t.Run("test caching with readEviction", func(t *testing.T) {
		err := readSubject.Write(ctx, "cache.test1", data)

		if err != nil {
			t.Error("writing to cache threw error", err)
		}

		second, err := readSubject.Read(ctx, "cache.test1")

		if err != nil {
			t.Error("reading from cache threw error", err)
		}

		if !bytes.Equal(data, second) {
			t.Error("data was not equal to the result", data, second)
		}

		<-time.After(time.Second)
		_, err = readSubject.Read(ctx, "cache.test1")

		if err != nil {
			if !errors.Is(err, lib.ErrKeyNotFound) {
				t.Error("ErrKeyNotFound was not returned as expected", err)
			}
		} else {
			t.Error("no error was thrown", data)
		}
	})

	t.Run("delete a cache item", func(t *testing.T) {
		err := writeSubject.Write(ctx, "cache.test1", data)

		if err != nil {
			t.Error("writing to cache threw error", err)
		}

		second, err := writeSubject.Read(ctx, "cache.test1")

		if err != nil {
			t.Error("reading from cache threw error", err)
		}

		if !bytes.Equal(data, second) {
			t.Error("data was not equal to the result", data, second)
		}

		err = writeSubject.Delete(ctx, "cache.test1")

		if err != nil {
			t.Error("deleting key from cache threw error", err)
		}

		_, err = writeSubject.Read(ctx, "cache.test1")

		if err != nil {
			if !errors.Is(err, lib.ErrKeyNotFound) {
				t.Error("ErrKeyNotFound was not returned as expected", err)
			}
		} else {
			t.Error("no error was thrown", data)
		}
	})
}
