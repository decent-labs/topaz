package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type App struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`

	Interval    int    `json:"interval"`
	Name        string `json:"name"`
	LastBatched *int64 `json:"-"`
	EthAddress  string `json:"ethAddress"`

	UserID string `json:"userId"`
	User   *User  `json:"user,omitempty"`
}

type Apps []App

func (a *App) CreateApp(db *gorm.DB) error {
	return db.Create(&a).Error
}

func (a *App) FindApp(db *gorm.DB) error {
	return db.Where(&a).First(&a).Error
}

func (as *Apps) GetAppsToBatch(db *gorm.DB) error {
	clause := "last_batched is null or (extract(epoch from now()) - last_batched >= interval)"
	return db.Where(clause).Find(&as).Error
}

func (a *App) UpdateApp(db *gorm.DB) error {
	return db.Save(a).Error
}

func (as *Apps) GetAppsForUser(u *User, db *gorm.DB) error {
	return db.Model(&u).Related(&as).Error
}
