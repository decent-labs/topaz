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

type Objects []Object

func (o *Object) CreateObject(db *gorm.DB) error {
	if err := db.Create(&o).Error; err != nil {
		return err
	}
	return nil
}

func (os *Objects) GetObjectsByAppID(db *gorm.DB, id uint) error {
	clause := "flush_id IS NULL AND app_id = ?"

	if err := db.Where(clause, id).Find(&os).Error; err != nil {
		return err
	}
	return nil
}
