package models

import (
	"strings"
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

	Hashes Hashes `json:"-"`
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

// GetProof ...
func (p *Proof) GetProof(db *gorm.DB) error {
	return db.Model(&p.App).Related(&p).Error
}

// CheckValidity ...
func (p *Proof) CheckValidity() error {
	cur, err := p.Hashes.GetMerkleRoot()
	if err != nil {
		return err
	}

	validRoot := strings.Compare(cur, p.MerkleRoot) == 0

	t, err := makeMerkleTree(&p.Hashes)
	if err != nil {
		return err
	}

	validTree, err := t.VerifyTree()
	if err != nil {
		return err
	}

	p.ValidStructure = validRoot && validTree
	return nil
}

