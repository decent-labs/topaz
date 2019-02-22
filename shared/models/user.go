package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User ...
type User struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Name     *string `json:"name,omitempty"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
}

// UpdatePassword ...
type UpdatePassword struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

// CreateUser ...
func (u *User) CreateUser(db *gorm.DB) error {
	return db.Create(&u).Error
}

// GetUser ...
func (u *User) GetUser(db *gorm.DB) error {
	return db.First(&u).Error
}

// UpdateUser ...
func (u *User) UpdateUser(db *gorm.DB) error {
	return db.Save(&u).Error
}

// GetUserWithEmail ...
func (u *User) GetUserWithEmail(db *gorm.DB) error {
	return db.Where(&User{Email: u.Email}).First(&u).Error
}
