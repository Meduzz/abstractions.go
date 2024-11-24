package caching_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/Meduzz/abstractions.go/internal/local/caching"
	"github.com/Meduzz/abstractions.go/lib"
)

func TestLocalCache(t *testing.T) {
	writeSubject := caching.NewCache(lib.EvictionWrite, 10*time.Millisecond)
	readSubject := caching.NewCache(lib.EvictionRead, 10*time.Millisecond)
	data1 := []byte("data1")
	ctx := context.Background()

	t.Run("test writeEviction", func(t *testing.T) {
		err := writeSubject.Write(ctx, "data1", data1)

		if err != nil {
			t.Errorf("error was not nil, %v", err)
		}

		result, err := writeSubject.Read(ctx, "data1")

		if err != nil {
			t.Errorf("error was not nil, %v", err)
		}

		if !bytes.Equal(result, data1) {
			t.Errorf("result was not data1 but %v", result)
		}
		timer := time.After(15 * time.Millisecond)

		<-timer

		result, err = writeSubject.Read(ctx, "data1")

		if err == nil {
			t.Error("error was nil")
		}

		if result != nil {
			t.Error("result was not nil")
		}
	})

	t.Run("test readEviction", func(t *testing.T) {
		err := readSubject.Write(ctx, "data1", data1)

		if err != nil {
			t.Errorf("error was not nil, %v", err)
		}

		result, err := readSubject.Read(ctx, "data1")

		if err != nil {
			t.Errorf("error was not nil, %v", err)
		}

		if !bytes.Equal(result, data1) {
			t.Errorf("result was not data1 but %v", result)
		}
		timer := time.After(15 * time.Millisecond)

		<-timer

		result, err = readSubject.Read(ctx, "data1")

		if err == nil {
			t.Error("error was nil")
		}

		if result != nil {
			t.Error("result was not nil")
		}
	})

	t.Run("remove an item from cache", func(t *testing.T) {
		err := writeSubject.Write(ctx, "data1", data1)

		if err != nil {
			t.Errorf("error was not nil, %v", err)
		}

		result, err := writeSubject.Read(ctx, "data1")

		if err != nil {
			t.Errorf("error was not nil, %v", err)
		}

		if !bytes.Equal(result, data1) {
			t.Errorf("result was not data1 but %v", result)
		}

		writeSubject.Delete(ctx, "data1")

		_, err = writeSubject.Read(ctx, "data1")

		if err == nil {
			t.Errorf("there was no error")
		}
	})
}
