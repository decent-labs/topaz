package models

import "github.com/jinzhu/gorm"

// Proof represents a collection of objects prepared for IPFS and Ethereum
type Proof struct {
	gorm.Model
	BatchID        uint
	Batch          Batch
	DirectoryHash  string
	EthTransaction string
}

// CreateProof creates a new entry in our database for a 'batch'
func (p *Proof) CreateProof(db *gorm.DB) error {
	return db.Create(&p).Error
}
