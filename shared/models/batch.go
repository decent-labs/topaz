package models

import "github.com/jinzhu/gorm"

// Batch represents a collection of objects prepared for IPFS and Ethereum
type Batch struct {
	gorm.Model
	UnixTimestamp int64
	AppID         uint
	App           App
	Objects       []Object
}

// CreateBatch creates a new entry in our database for a 'batch'
func (b *Batch) CreateBatch(db *gorm.DB) error {
	if err := db.Create(&b).Error; err != nil {
		return err
	}
	return nil
}
