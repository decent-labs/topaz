package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// App ...
type App struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

	Interval    int    `json:"interval"`
	Name        string `json:"name"`
	LastBatched *int64 `json:"-"`
	EthAddress  string `json:"ethAddress"`

	UserID string `json:"userId"`
	User   *User  `json:"-"`
}

// Apps ...
type Apps []App

// CreateApp ...
func (a *App) CreateApp(db *gorm.DB) error {
	return db.Create(&a).Error
}

// GetApps ...
func (as *Apps) GetApps(a *App, db *gorm.DB) error {
	return db.Model(&a.User).Related(&as).Error
}

// GetApp ...
func (a *App) GetApp(db *gorm.DB) error {
	return db.Model(&a.User).Related(&a).Error
}

// For Batching

// GetAppsToBatch ...
func (as *Apps) GetAppsToBatch(db *gorm.DB) error {
	clause := "last_batched is null or (extract(epoch from now()) - last_batched >= interval)"
	return db.Where(clause).Find(&as).Error
}

// UpdateApp ...
func (a *App) UpdateApp(db *gorm.DB) error {
	return db.Save(a).Error
}
