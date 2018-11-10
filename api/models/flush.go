package models

import "github.com/jinzhu/gorm"

type Flush struct {
	gorm.Model
	DirectoryHash string
	Transaction   string
	AppID         int
	App           App
	Objects       []Object
}
