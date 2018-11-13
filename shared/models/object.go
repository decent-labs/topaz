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
	return db.Create(&o).Error
}

func (os *Objects) GetObjectsByAppID(db *gorm.DB, id uint) error {
	clause := "flush_id IS NULL AND app_id = ?"
	return db.Where(clause, id).Find(&os).Error
}

func (os *Objects) GetObjectsByHash(db *gorm.DB, o *Object) error {
	return db.Where(&Object{Hash: o.Hash, AppID: o.AppID}).Find(&os).Error
}