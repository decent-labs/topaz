package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) CreateUser(db *gorm.DB) error {
	if err := db.Create(&u).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) GetUser(db *gorm.DB) error {
	if err := db.Where("email = ?", u.Email).First(&u).Error; err != nil {
		return err
	}
	return nil
}
