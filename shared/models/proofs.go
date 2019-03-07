package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Proof ...
type Proof struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	MerkleRoot     string `json:"merkleRoot"`
	EthTransaction string `json:"ethTransaction"`
	UnixTimestamp  int64  `json:"unixTimestamp"`

	AppID string `json:"appId"`
	App   *App   `json:"-"`

	HashStubs HashStubs `json:"hashes,omitempty"`
}

// Proofs ...
type Proofs []Proof

// CreateProof ...
func (p *Proof) CreateProof(db *gorm.DB) error {
	return db.Create(&p).Error
}

// GetProofs ...
func (ps *Proofs) GetProofs(p *Proof, db *gorm.DB) error {
	return db.Model(&p.App).Related(&ps).Error
}

// GetProofWithHashStubs ...
func (p *Proof) GetProofWithHashStubs(db *gorm.DB) error {
	if err := db.Model(&p.App).Related(&p).Error; err != nil {
		return err
	}

	hs := HashStubs{}
	if err := hs.GetHashesByProof(db, p); err != nil {
		return err
	}

	p.HashStubs = hs
	return nil
}
