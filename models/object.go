package models

import "github.com/jinzhu/gorm"

type Object struct {
	gorm.Model
	DataBlob []byte `json:"dataBlob"`
	Hash     string `json:"hash"`
	AppID    int    `json:"appId"`
	App      App    `json:"app"`
	FlushID  *int   `json:"flushId"`
	Flush    Flush  `json:"flush"`
}

func (o *Object) CreateObject(db *gorm.DB) error {
	if err := db.Create(&o).Error; err != nil {
		return err
	}
	return nil
}
