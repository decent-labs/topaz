package models

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	b58 "github.com/mr-tron/base58/base58"
)

// APIToken ...
type APIToken struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Token string `json:"token"`

	UserID string `json:"userId"`
	User   *User  `json:"-"`
}

// APITokens ...
type APITokens []APIToken

// CreateAPIToken ...
func (a *APIToken) CreateAPIToken(db *gorm.DB) error {
	if err := a.makeNewToken(); err != nil {
		return err
	}
	return db.Create(&a).Error
}

// GetAPITokens ...
func (as *APITokens) GetAPITokens(a *APIToken, db *gorm.DB) error {
	return db.Model(&a.User).Related(&as).Error
}

// GetAPIToken ...
func (a *APIToken) GetAPIToken(db *gorm.DB) error {
	return db.Model(&a.User).Related(&a).Error
}

// GetAPITokenFromToken ...
func (a *APIToken) GetAPITokenFromToken(db *gorm.DB) error {
	if a.Token == "" {
		return errors.New("no token")
	}
	return db.Where(&APIToken{Token: a.Token}).First(&a).Error
}

func (a *APIToken) makeNewToken() error {
	rb := make([]byte, 256)
	_, err := rand.Read(rb)
	if err != nil {
		return err
	}
	hash := sha256.Sum256(rb)
	token := b58.Encode(hash[:])
	a.Token = token
	return nil
}
