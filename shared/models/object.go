package models

import "github.com/jinzhu/gorm"

type Object struct {
	gorm.Model
	DataBlob []byte `json:"dataBlob"`
	Hash     string `json:"hash"`
	AppID    uint   `json:"appId"`
	App      App    `json:"app"`
	BatchID  *uint  `json:"batchId"`
	Batch    Batch  `json:"batch"`
}

type Objects []Object

func (o *Object) CreateObject(db *gorm.DB) error {
	return db.Create(&o).Error
}

func (os *Objects) GetObjectsByAppID(db *gorm.DB, id uint) error {
	clause := "batch_id IS NULL AND app_id = ?"
	return db.Where(clause, id).Find(&os).Error
}

func (os *Objects) GetObjectsByHash(db *gorm.DB, o *Object) error {
	return db.Where(&Object{Hash: o.Hash, AppID: o.AppID}).Find(&os).Error
}

func (o *Object) UpdateObject(db *gorm.DB) error {
	if err := db.Save(&o).Error; err != nil {
		return err
	}
	return nil
}
