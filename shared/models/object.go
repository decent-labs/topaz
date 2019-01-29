package models

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

type Object struct {
	gorm.Model

	UUID          string `json:"uuid"`
	UnixTimestamp int64  `json:"unixTimestamp"`

	AppID uint `json:"appId"`
	App   *App `json:"app,omitempty"`
}

type Objects []Object

func (o *Object) CreateObject(db *gorm.DB) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	o.UUID = uuid.String()
	return db.Create(&o).Error
}

func (os *Objects) GetObjectsByTimestamps(db *gorm.DB, appId uint, start int, end int) error {
	return db.
		Preload("Proof.Batch").
		Where("app_id = ?", appId).
		Where("unix_timestamp BETWEEN (?) AND (?)", start, end).
		Find(&os).
		Error
}
