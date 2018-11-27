package models

import "github.com/jinzhu/gorm"

type Batch struct {
	gorm.Model

	UnixTimestamp int64 `json:"unixTimestamp"`

	AppID uint `json:"appId"`
	App   *App `json:"app,omitempty"`

	Objects *Objects `json:"objects,omitempty"`
}

type Batches []Batch

func (b *Batch) CreateBatch(db *gorm.DB) error {
	if err := db.Create(&b).Error; err != nil {
		return err
	}
	return nil
}
