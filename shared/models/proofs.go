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

	MerkleRoot    string `json:"merkleRoot"`
	UnixTimestamp int64  `json:"unixTimestamp"`

	AppID string `json:"appId"`
	App   *App   `json:"-"`

	HashStubs    HashStubs              `json:"hashes,omitempty"`
	Transactions BlockchainTransactions `json:"blockchainTransactions,omitempty"`
}

// Proofs ...
type Proofs []Proof

// GetProofs ...
func (ps *Proofs) GetProofs(p *Proof, db *gorm.DB) error {
	return db.Model(&p.App).Order("created_at").Related(&ps).Error
}

// GetFullProof ...
func (p *Proof) GetFullProof(db *gorm.DB) error {
	if err := db.Model(&p.App).Related(&p).Error; err != nil {
		return err
	}

	hs := HashStubs{}
	if err := hs.GetHashesByProof(db, p); err != nil {
		return err
	}
	p.HashStubs = hs

	bts := BlockchainTransactions{}
	if err := bts.GetBlockchainTransactionsByProof(db, p); err != nil {
		return err
	}

	for i, bt := range bts {
		bn := BlockchainNetwork{ID: bt.BlockchainNetworkID}
		if err := bn.GetBlockchainNetwork(db); err != nil {
			return err
		}
		bts[i].BlockchainNetworkName = bn.Name

		bes := BlockchainExplorers{}
		if err := bes.GetBlockchainExplorersByNetworkID(db, bt.BlockchainNetworkID); err != nil {
			return err
		}

		var urls []string
		for _, be := range bes {
			s := strings.Replace(be.URLTemplate, "{transaction_hash}", bt.TransactionHash, 1)
			urls = append(urls, s)
		}

		bts[i].Explorers = urls
	}

	p.Transactions = bts

	return nil
}
