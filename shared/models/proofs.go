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
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

	MerkleRoot     string `json:"merkleRoot"`
	EthTransaction string `json:"ethTransaction"`
	ValidStructure bool   `json:"validStructure" gorm:"-"`

	BatchID string `json:"batchId"`
	Batch   *Batch `json:"batch,omitempty"`

	Hashes Hashes `json:"-"`
}

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
