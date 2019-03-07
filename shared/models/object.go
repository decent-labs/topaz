package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Object ...
type Object struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	AppID string `json:"appId"`
	App   *App   `json:"-"`

	HashStubs HashStubs `json:"hashes,omitempty"`
}

// Objects ...
type Objects []Object

// CreateObject ...
func (o *Object) CreateObject(db *gorm.DB) error {
	return db.Create(&o).Error
}

// GetObjects ...
func (os *Objects) GetObjects(o *Object, db *gorm.DB) error {
	return db.Model(&o.App).Order("created_at").Related(&os).Error
}

// GetObject ...
func (o *Object) GetObject(db *gorm.DB) error {
	return db.Model(&o.App).Related(&o).Error
}

// GetObjectWithHashStubs ...
func (o *Object) GetObjectWithHashStubs(db *gorm.DB) error {
	if err := db.Model(&o.App).Related(&o).Error; err != nil {
		return err
	}

	hs := HashStubs{}
	if err := hs.GetHashesByObject(db, o); err != nil {
		return err
	}

	o.HashStubs = hs
	return nil
}
