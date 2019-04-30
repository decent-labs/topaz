package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// App ...
type App struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Interval    int    `json:"interval"`
	Name        string `json:"name"`
	LastProofed *int64 `json:"-"`

	UserID string `json:"userId"`
	User   *User  `json:"-"`

	Proofs Proofs `json:"-"`
}

// Apps ...
type Apps []App

// CreateApp ...
func (a *App) CreateApp(db *gorm.DB) error {
	return db.Create(&a).Error
}

// GetApps ...
func (as *Apps) GetApps(a *App, db *gorm.DB) error {
	return db.Model(&a.User).Order("created_at").Related(&as).Error
}

// GetApp ...
func (a *App) GetApp(db *gorm.DB) error {
	return db.Model(&a.User).Related(&a).Error
}
