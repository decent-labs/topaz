package models

import (
	"strings"

	"github.com/jinzhu/gorm"
)

type Proof struct {
	gorm.Model

	MerkleRoot     string `json:"merkleRoot"`
	EthTransaction string `json:"ethTransaction"`
	ValidStructure bool   `json:"validStructure" gorm:"-"`

	BatchID uint   `json:"batchId"`
	Batch   *Batch `json:"batch,omitempty"`

	Hashes Hashes `json:"hashes,omitempty"`
}

func (p *Proof) CreateProof(db *gorm.DB) error {
	return db.Create(&p).Error
}

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
