package models

import "github.com/jinzhu/gorm"

type Object struct {
	gorm.Model
	Hash    string
	AppID   int
	App     App
	FlushID *int
	Flush   Flush
}

type StoreResponse struct {
	Hash  string `json:"hash"`
	AppID string `json:"appId"`
}
