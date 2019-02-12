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
	EthAddress  string `json:"ethAddress"`

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
	return db.Model(&a.User).Related(&as).Error
}

// GetApp ...
func (a *App) GetApp(db *gorm.DB) error {
	return db.Model(&a.User).Related(&a).Error
}

// For proofing

// GetAppsToProof ...
func (as *Apps) GetAppsToProof(db *gorm.DB) error {
	clause := "last_proofed is null or (extract(epoch from now()) - last_proofed >= interval)"
	return db.Where(clause).Find(&as).Error
}

// UpdateApp ...
func (a *App) UpdateApp(db *gorm.DB) error {
	return db.Save(a).Error
}
