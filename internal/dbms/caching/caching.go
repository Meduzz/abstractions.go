package caching

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Meduzz/abstractions.go/internal/interval"
	"github.com/Meduzz/abstractions.go/lib"
	"gorm.io/gorm"
)

type (
	CacheItem struct {
		ID      int64     `gorm:"primaryKey,autoIncrement"`
		Key     string    `gorm:"size:256"`
		Value   []byte    `gorm:"type:bytes"`
		Created time.Time `gorm:"autoCreateTime:milli"`
	}

	delegate struct {
		db       *gorm.DB
		prefix   string
		eviction lib.Eviction
		ttl      time.Duration
	}
)

func NewDbmsCachingDelegate(db *gorm.DB, eviction lib.Eviction, prefix string, ttl time.Duration) (lib.CacheStorageDelegate, error) {
	err := db.AutoMigrate(&CacheItem{})

	if err != nil {
		return nil, err
	}

	interval.OnInterval(time.Minute, func() {
		age := time.Now().Add(-1 * ttl)
		db.Where("created < ?", age).Delete(&CacheItem{})
	})

	return &delegate{
		db:       db,
		prefix:   prefix,
		eviction: eviction,
		ttl:      ttl,
	}, nil
}

func (d *delegate) Write(ctx context.Context, key string, data []byte) error {
	item := &CacheItem{}
	item.Key = d.withPrefix(key)
	item.Value = data

	return d.db.Save(item).Error
}

func (d *delegate) Read(ctx context.Context, key string) ([]byte, error) {
	item := &CacheItem{}

	err := d.db.Where("key = ?", d.withPrefix(key)).Last(item).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, lib.ErrKeyNotFound
		}

		return nil, err
	}

	if item.Created.Add(d.ttl).Before(time.Now()) {
		d.db.Delete(item)
		return nil, lib.ErrKeyNotFound
	}

	if d.eviction == lib.EvictionRead {
		item.Created = time.Now()
		d.db.Save(item)
	}

	return item.Value, nil
}

func (d *delegate) Delete(ctx context.Context, key string) error {
	return d.db.Where("key = ?", d.withPrefix(key)).Delete(&CacheItem{}).Error
}

func (d *delegate) withPrefix(key string) string {
	return fmt.Sprintf("%s.%s", d.prefix, key)
}
