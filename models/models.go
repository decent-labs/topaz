package models

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// Auth structures

type AuthAdminClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// API Request structures

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateAdminTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateAppRequest struct {
	Interval int    `json:"interval"`
	Name     string `json:"name"`
}

type CreateAppTokenRequest struct {
	AppId int `json:"appId"`
}

// API Response structures

type Exception struct {
	Message string `json:"message"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type DeployResponse struct {
	Tx   string
	Addr string
}

type StoreResponse struct {
	Hash  string `json:"hash"`
	AppID string `json:"appId"`
}

// Database models

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

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
