package models

import (
	"encoding/json"
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

// MarshalJSON ...
func (p *Proof) MarshalJSON() ([]byte, error) {
	type Alias Proof
	return json.Marshal(&struct {
		*Alias
		Hashes []interface{} `json:"hashes"`
	}{
		Alias:  (*Alias)(p),
		Hashes: p.reduceHashes(),
	})
}

// CreateProof ...
func (p *Proof) CreateProof(db *gorm.DB) error {
	return db.Create(&p).Error
}

// GetProofs ...
func (ps *Proofs) GetProofs(p *Proof, db *gorm.DB) error {
	return db.Model(&p.Batch).Related(&ps).Error
}

// GetProof ...
func (p *Proof) GetProof(db *gorm.DB) error {
	return db.Model(&p.Batch).Related(&p).Error
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

func (p *Proof) reduceHashes() []interface{} {
	a := make([]interface{}, 0, len(p.Hashes))
	for _, h := range p.Hashes {
		a = append(a, struct {
			ID   string `json:"id"`
			Hash string `json:"hash"`
		}{
			h.ID,
			h.TransformHashToHex(),
		})
	}
	return a
}
