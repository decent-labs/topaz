package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type App struct {
	gorm.Model
	Interval    int        `json:"interval"`
	Name        string     `json:"name"`
	LastBatched *time.Time `json:"lastBatched"`
	UserID      uint       `json:"userID"`
	User        User       `json:"user"`
	Batches     []Batch    `json:"batches"`
	Objects     []Object   `json:"objects"`
	EthAddress  string     `json:"ethAddress"`
}

type Apps []App

func (a *App) CreateApp(db *gorm.DB) error {
	return db.Create(&a).Error
}

func (a *App) GetApp(db *gorm.DB) error {
	return db.Where("id = ? AND user_id = ?", a.ID, a.UserID).First(&a).Error
}

func (as *Apps) GetAppsToBatch(db *gorm.DB) error {
	// TODO: Fix query
	// clause := "last_batched IS NULL OR NOW() - last_batched >= interval * '1 second'::interval"

	return db.Find(&as).Error
}

func (a *App) UpdateApp(db *gorm.DB) error {
	return db.Save(a).Error
}
