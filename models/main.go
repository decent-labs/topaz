package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

type App struct {
	gorm.Model
	Interval    int
	Name        string
	LastFlushed *time.Time
	UserID      int
	User        User
	Flushes     []Flush
	Objects     []Object
}

type Flush struct {
	gorm.Model
	DirectoryHash string
	Transaction   string
	AppID         int
	App           App
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
