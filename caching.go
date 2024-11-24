package abstractions

import (
	"context"

	"github.com/Meduzz/abstractions.go/lib"
)

type (
	Serializer   = func(any) ([]byte, error)
	Deserializer = func([]byte) (any, error)

	Caching struct {
		storage      lib.CacheStorageDelegate
		serializer   Serializer
		deserializer Deserializer
	}
)

func NewCaching(storage lib.CacheStorageDelegate, serializer Serializer, deserializer Deserializer) *Caching {
	return &Caching{
		storage:      storage,
		serializer:   serializer,
		deserializer: deserializer,
	}
}

func (c *Caching) Write(ctx context.Context, key string, data any) error {
	bs, err := c.serializer(data)

	if err != nil {
		return err
	}

	return c.storage.Write(ctx, key, bs)
}

func (c *Caching) Read(ctx context.Context, key string) (any, error) {
	bs, err := c.storage.Read(ctx, key)

	if err != nil {
		return nil, err
	}

	return c.deserializer(bs)
}

func (c *Caching) Delete(ctx context.Context, key string) error {
	return c.storage.Delete(ctx, key)
}
