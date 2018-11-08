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
