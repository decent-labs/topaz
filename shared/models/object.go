package models

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

type Object struct {
	gorm.Model

	UUID string `json:"uuid"`

	AppID uint `json:"appId"`
	App   *App `json:"app,omitempty"`
}

type Objects []Object

func (o *Object) MakeUUID() error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	o.UUID = uuid.String()
	return nil
}

func (o *Object) CreateObject(db *gorm.DB) error {
	return db.Create(&o).Error
}

func (o *Object) FindObject(db *gorm.DB) error {
	return db.Where(o).First(&o).Error
}
