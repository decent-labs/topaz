package models

import "github.com/jinzhu/gorm"

type Proof struct {
	gorm.Model

	DirectoryHash  string `json:"directoryHash"`
	EthTransaction string `json:"ethTransaction"`

	BatchID uint   `json:"batchId"`
	Batch   *Batch `json:"batch,omitempty"`
}

func (p *Proof) CreateProof(db *gorm.DB) error {
	return db.Create(&p).Error
}
