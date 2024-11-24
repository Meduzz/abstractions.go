package caching

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Meduzz/abstractions.go/lib"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDbmsCaching(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	writeSubject, err := NewDbmsCachingDelegate(db, lib.EvictionWrite, "write", time.Second)

	if err != nil {
		t.Error(err)
	}

	readSubject, err := NewDbmsCachingDelegate(db, lib.EvictionRead, "read", time.Second)

	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	data := []byte("Im data")

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
