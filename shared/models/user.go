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
	DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

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
	return db.First(&u).Error
}

// GetUserWithEmail ...
func (u *User) GetUserWithEmail(db *gorm.DB) error {
	return db.Where(&User{Email: u.Email}).First(&u).Error
}
