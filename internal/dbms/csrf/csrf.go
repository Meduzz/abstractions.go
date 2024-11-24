package csrf

import (
	"context"
	"errors"
	"time"

	"github.com/Meduzz/abstractions.go/internal/interval"
	"github.com/Meduzz/abstractions.go/lib"
	"gorm.io/gorm"
)

type (
	abstraction struct {
		db  *gorm.DB
		ttl time.Duration
	}

	CsrfEntity struct {
		ID      int64     `gorm:"primaryKey,autoIncrement"`
		Created time.Time `gorm:"autoCreateTime:milli"`
		Key     string    `gorm:"size:64,index"`
		Value   string    `gorm:"size:128"`
	}
)

func NewDbmsCSRFStorageDelegate(db *gorm.DB, ttl time.Duration) (lib.CSRFStorageDelegate, error) {
	err := db.AutoMigrate(&CsrfEntity{})

	if err != nil {
		return nil, err
	}

	interval.OnInterval(time.Minute, func() {
		age := time.Now().Add(-1 * ttl)
		db.Where("created < ?", age).Delete(&CsrfEntity{})
	})

	return &abstraction{db, ttl}, nil
}

func (d *abstraction) Store(ctx context.Context, token *lib.CSRFToken) error {
	item := &CsrfEntity{}
	item.Key = token.Key
	item.Value = token.Value

	return d.db.Save(item).Error
}

func (d *abstraction) Verify(ctx context.Context, token *lib.CSRFToken) (bool, error) {
	item := &CsrfEntity{}
	err := d.db.Where("key = ?", token.Key).First(item).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	valid := item.Value == token.Value && item.Created.Add(d.ttl).After(time.Now())

	err = d.db.Delete(item).Error

	if err != nil {
		return valid, err
	}

	return valid, nil
}
