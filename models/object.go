package models

import "github.com/jinzhu/gorm"

type Object struct {
	gorm.Model
	DataBlob []byte `json:"dataBlob"`
	Hash     string `json:"hash"`
	AppID    int    `json:"appId"`
	App      App    `json:"app"`
	FlushID  *int   `json:"flushId"`
	Flush    Flush  `json:"flush"`
}
