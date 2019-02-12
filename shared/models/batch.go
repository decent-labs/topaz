package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Batch ...
type Batch struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

	UnixTimestamp int64 `json:"unixTimestamp"`

	AppID string `json:"appId"`
	App   *App   `json:"-"`
}

// Batches ...
type Batches []Batch

// CreateBatch ...
func (b *Batch) CreateBatch(db *gorm.DB) error {
	return db.Create(&b).Error
}

// GetBatches ...
func (bs *Batches) GetBatches(b *Batch, db *gorm.DB) error {
	return db.Model(&b.App).Related(&bs).Error
}

// GetBatch ...
func (b *Batch) GetBatch(db *gorm.DB) error {
	return db.Model(&b.App).Related(&b).Error
}
