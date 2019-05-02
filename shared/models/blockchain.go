package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// BlockchainNetwork ...
type BlockchainNetwork struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Name string `json:"name"`
}

// BlockchainExplorer ...
type BlockchainExplorer struct {
	ID        string     `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	BlockchainNetworkID string             `json:"-"`
	BlockchainNetwork   *BlockchainNetwork `json:"-"`

	URLTemplate string `json:"urlTemplate"`
}

// BlockchainExplorers ...
type BlockchainExplorers []BlockchainExplorer

// BlockchainTransaction ...
type BlockchainTransaction struct {
	ID        string     `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	ProofID string `json:"-"`
	Proof   *Proof `json:"-"`

	BlockchainNetworkID   string             `json:"-"`
	BlockchainNetwork     *BlockchainNetwork `json:"-"`
	BlockchainNetworkName string             `json:"blockchainNetwork" gorm:"-"`

	TransactionHash string `json:"transactionHash"`

	Explorers []string `json:"explorers" gorm:"-"`
}

// BlockchainTransactions ...
type BlockchainTransactions []BlockchainTransaction

// CreateBlockchainTransaction ...
func (bt *BlockchainTransaction) CreateBlockchainTransaction(db *gorm.DB) error {
	return db.Create(&bt).Error
}

// GetBlockchainNetwork ...
func (bn *BlockchainNetwork) GetBlockchainNetwork(db *gorm.DB) error {
	return db.Find(&bn).Error
}

// GetBlockchainNetworkFromName ...
func (bn *BlockchainNetwork) GetBlockchainNetworkFromName(db *gorm.DB, name string) error {
	return db.First(&bn, &BlockchainNetwork{Name: name}).Error
}

// GetBlockchainTransactionsByProof ..
func (bts *BlockchainTransactions) GetBlockchainTransactionsByProof(db *gorm.DB, p *Proof) error {
	return db.
		Table("blockchain_transactions").
		Where(&BlockchainTransaction{ProofID: p.ID}).
		Order("created_at").
		Find(&bts).
		Error
}
