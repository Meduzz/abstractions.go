package log

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Meduzz/abstractions.go/lib"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	dbmsLogDelegate struct {
		db *gorm.DB
	}

	WorkItemStorage struct {
		ID   int64  `gorm:"primaryKey,autoIncrement"`
		Kind string `gorm:"size:128"`
		Work datatypes.JSON
	}
)

func NewDbmsWorkLogDelegate(db *gorm.DB) (lib.WorkLogDelegate, error) {
	err := db.AutoMigrate(&WorkItemStorage{})

	if err != nil {
		return nil, err
	}

	return &dbmsLogDelegate{db}, nil
}

// Append - append work to the queue
func (d *dbmsLogDelegate) Append(ctx context.Context, workItem *lib.WorkItem) error {
	log := &WorkItemStorage{}
	log.Kind = workItem.Kind
	log.Work = datatypes.JSON(workItem.Work)

	return d.db.Save(log).Error
}

// Size - fetch the size of the queue
func (d *dbmsLogDelegate) Size(ctx context.Context) (int64, error) {
	count := int64(0)
	err := d.db.Model(&WorkItemStorage{}).Count(&count).Error

	return count, err
}

// Fetch - fetch the first item in the log and remove it
func (d *dbmsLogDelegate) Fetch(ctx context.Context) (*lib.WorkItem, error) {
	work := &WorkItemStorage{}
	err := d.db.First(work).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	workItem := &lib.WorkItem{}
	workItem.Kind = work.Kind
	workItem.Work = json.RawMessage(work.Work)

	affected := d.db.Delete(work).RowsAffected

	if affected == 0 {
		return d.Fetch(ctx)
	}

	return workItem, nil
}
