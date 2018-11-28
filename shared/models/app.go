package models

import "github.com/jinzhu/gorm"

type App struct {
	gorm.Model
	Interval    int      `json:"interval"`
	Name        string   `json:"name"`
	LastBatched *int64   `json:"lastBatched"`
	UserID      uint     `json:"userID"`
	User        User     `json:"user"`
	Batches     []Batch  `json:"batches"`
	Objects     []Object `json:"objects"`
	EthAddress  string   `json:"ethAddress"`
}

type Apps []App

func (a *App) CreateApp(db *gorm.DB) error {
	return db.Create(&a).Error
}

func (a *App) GetApp(db *gorm.DB) error {
	return db.Model(&a).Related(&a.User).First(&a).Error
}

func (as *Apps) GetAppsToBatch(db *gorm.DB) error {
	clause := "last_batched is null or (extract(epoch from now()) - last_batched >= interval)"
	return db.Where(clause).Find(&as).Error
}

func (a *App) UpdateApp(db *gorm.DB) error {
	return db.Save(a).Error
}
