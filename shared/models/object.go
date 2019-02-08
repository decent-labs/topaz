package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Object struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`

	AppID string `json:"appId"`
	App   *App   `json:"app,omitempty"`

	Hashes Hashes `json:"-"`
}

type Objects []Object

func (o *Object) CreateObject(db *gorm.DB) error {
	return db.Create(&o).Error
}

func (o *Object) FindObject(db *gorm.DB) error {
	return db.Where(o).First(&o).Error
}

func (o *Object) FindFullObject(db *gorm.DB) error {
	return db.Preload("Proof.Hashes").Model(&o).Related(&o.Hashes).Error
}
