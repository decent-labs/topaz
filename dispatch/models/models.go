package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name string
}

type App struct {
	gorm.Model
	Name        string
	UserID      int
	User        User
	Interval    string
	LastFlushed *time.Time
	Flushes     []Flush
	Objects     []Object
}

type Flush struct {
	gorm.Model
	Transaction   string
	AppID         int
	App           App
	DirectoryHash string
	Objects       []Object
}

type Object struct {
	gorm.Model
	Hash    string
	AppID   int
	App     App
	FlushID *int
	Flush   Flush
}
