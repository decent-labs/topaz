package models

import (
	"encoding/json"
	"time"

	"github.com/decentorganization/topaz/api/crypto"
	"github.com/jinzhu/gorm"
)

// HashStub ...
type HashStub struct {
	ID      string `json:"id"`
	HashHex string `json:"hash"`
	Hash    []byte `json:"-"`
}

// HashStubs ...
type HashStubs []HashStub

// Hash ...
type Hash struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	HashHex       string `json:"hash" gorm:"-"`
	Hash          []byte `json:"-"`
	UnixTimestamp int64  `json:"unixTimestamp"`

	ObjectID *string `json:"objectId"`
	Object   *Object `json:"-"`
	ProofID  *string `json:"proofId"`
	Proof    *Proof  `json:"-"`
}

// Hashes ...
type Hashes []Hash

// MarshalJSON ...
func (h *Hash) MarshalJSON() ([]byte, error) {
	type Alias Hash
	return json.Marshal(&struct {
		*Alias
		HashHex string `json:"hash"`
	}{
		Alias:   (*Alias)(h),
		HashHex: crypto.TransformHashToHex(h.Hash),
	})
}

// MarshalJSON ...
func (hs *HashStub) MarshalJSON() ([]byte, error) {
	type Alias HashStub
	return json.Marshal(&struct {
		*Alias
		HashHex string `json:"hash"`
	}{
		Alias:   (*Alias)(hs),
		HashHex: crypto.TransformHashToHex(hs.Hash),
	})
}

// CreateHash ...
func (h *Hash) CreateHash(db *gorm.DB) error {
	return db.Create(&h).Error
}

// GetHashes ...
func (hs *Hashes) GetHashes(h *Hash, db *gorm.DB) error {
	return db.Model(&h.Object).Related(&hs).Error
}

// GetHash ...
func (h *Hash) GetHash(db *gorm.DB) error {
	return db.Model(&h.Object).Related(&h).Error
}

// MakeMerkleLeafs ...
func (hs *Hashes) MakeMerkleLeafs() crypto.MerkleLeafs {
	var ms crypto.MerkleLeafs
	for _, h := range *hs {
		m := crypto.MerkleLeaf{Hash: h.Hash}
		ms = append(ms, m)
	}
	return ms
}

// GetHashesByApp ...
func (hs *Hashes) GetHashesByApp(db *gorm.DB, app *App) error {
	return db.
		Table("hashes").
		Joins("join objects on objects.id = hashes.object_id").
		Joins("join apps on apps.id = objects.app_id").
		Where("apps.id = (?)", app.ID).
		Where("hashes.proof_id IS NULL").
		Find(&hs).
		Error
}

// GetHashesByProof ...
func (hs *HashStubs) GetHashesByProof(db *gorm.DB, p *Proof) error {
	return db.
		Table("hashes").
		Where(&Hash{ProofID: &p.ID}).
		Find(&hs).
		Error
}

// UpdateWithProof ...
func (hs Hashes) UpdateWithProof(db *gorm.DB, proofID *string) error {
	ids := make([]string, len(hs))
	for i, h := range hs {
		ids[i] = h.ID
	}

	return db.Model(Hash{}).Where("id IN (?)", ids).Updates(Hash{ProofID: proofID}).Error
}
