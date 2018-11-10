package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type App struct {
	gorm.Model
	Interval    int        `json:"interval"`
	Name        string     `json:"name"`
	LastFlushed *time.Time `json:"lastFlushed"`
	UserID      int        `json:"userID"`
	User        User       `json:"user"`
	Flushes     []Flush    `json:"flushes"`
	Objects     []Object   `json:"objects"`
	EthAddress  string     `json:"ethAddress"`
}

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
