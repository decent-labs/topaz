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
	if err := db.Create(&a).Error; err != nil {
		return err
	}
	return nil
}

func (a *App) GetApp(db *gorm.DB) error {
	if err := db.Where("id = ? AND user_id = ?", a.ID, a.UserID).First(&a).Error; err != nil {
		return err
	}
	return nil
}

func (as *Apps) GetAppsToBatch(db *gorm.DB) error {
	clause := "last_batched IS NULL OR NOW() - last_batched >= interval * '1 second'::interval"

	if err := db.Where(clause).Find(&as).Error; err != nil {
		return err
	}
	return nil
}
