package models

import "github.com/jinzhu/gorm"

type Flush struct {
	gorm.Model
	DirectoryHash  string
	EthTransaction string
	AppID          uint
	App            App
	Objects        []Object
}

func (f *Flush) CreateFlush(db *gorm.DB) error {
	if err := db.Create(&f).Error; err != nil {
		return err
	}
	return nil
}
