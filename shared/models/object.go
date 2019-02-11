package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Object ...
type Object struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

	AppID string `json:"appId"`
	App   *App   `json:"-"`

	Hashes Hashes `json:"-"`
}

// Objects ...
type Objects []Object

// CreateObject ...
func (o *Object) CreateObject(db *gorm.DB) error {
	return db.Create(&o).Error
}

// GetObject ...
func (o *Object) GetObject(db *gorm.DB) error {
	return db.Model(&o.App).Related(&o).Error
}

// GetObjects ...
func (os *Objects) GetObjects(o *Object, db *gorm.DB) error {
	return db.Model(&o.App).Related(&os).Error
}
