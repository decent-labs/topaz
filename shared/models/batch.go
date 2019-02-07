package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Batch struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`

	UnixTimestamp int64 `json:"unixTimestamp"`

	AppID string `json:"appId"`
	App   *App   `json:"app,omitempty"`
}

type Batches []Batch

func (b *Batch) CreateBatch(db *gorm.DB) error {
	if err := db.Create(&b).Error; err != nil {
		return err
	}
	return nil
}
