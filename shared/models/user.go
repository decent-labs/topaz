package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User ...
type User struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`

	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUser ...
func (u *User) CreateUser(db *gorm.DB) error {
	return db.Create(&u).Error
}

// GetUser ...
func (u *User) GetUser(db *gorm.DB) error {
	if err := db.Where("email = ?", u.Email).First(&u).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) FindUser(db *gorm.DB) error {
	return db.First(&u).Error
}
