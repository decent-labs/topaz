package models

import (
	"encoding/json"
	"time"

	"github.com/decentorganization/topaz/api/crypto"
	"github.com/jinzhu/gorm"
	multihash "github.com/multiformats/go-multihash"
)

// HashStub ...
type HashStub struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	HashString string `json:"hash" gorm:"-"`
	MultiHash  []byte `json:"-" gorm:"column:hash"`
}

// HashStubs ...
type HashStubs []HashStub

// Hash ...
type Hash struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	HashString    string `json:"hash" gorm:"-"`
	MultiHash     []byte `json:"-" gorm:"column:hash"`
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
		HashString string `json:"hash"`
	}{
		Alias:      (*Alias)(h),
		HashString: HashBytesToString(h.MultiHash),
	})
}

// MarshalJSON ...
func (hs *HashStub) MarshalJSON() ([]byte, error) {
	type Alias HashStub
	return json.Marshal(&struct {
		*Alias
		HashString string `json:"hash"`
	}{
		Alias:      (*Alias)(hs),
		HashString: HashBytesToString(hs.MultiHash),
	})
}

// HashBytesToString ...
func HashBytesToString(hash []byte) string {
	var mh multihash.Multihash = hash
	return mh.B58String()
}

// CreateHash ...
func (h *Hash) CreateHash(db *gorm.DB) error {
	return db.Create(&h).Error
}

// GetHashes ...
func (hs *Hashes) GetHashes(h *Hash, db *gorm.DB) error {
	return db.Model(&h.Object).Order("created_at").Related(&hs).Error
}

// GetHash ...
func (h *Hash) GetHash(db *gorm.DB) error {
	return db.Model(&h.Object).Related(&h).Error
}

// MakeMerkleLeafs ...
func (hs *Hashes) MakeMerkleLeafs() crypto.MerkleLeafs {
	var ms crypto.MerkleLeafs
	for _, h := range *hs {
		m := crypto.MerkleLeaf{Hash: h.MultiHash}
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
		Order("hashes.created_at").
		Find(&hs).
		Error
}

// GetHashesByProof ...
func (hs *HashStubs) GetHashesByProof(db *gorm.DB, p *Proof) error {
	return db.
		Table("hashes").
		Where(&Hash{ProofID: &p.ID}).
		Order("created_at").
		Find(&hs).
		Error
}

// GetHashesByObject ...
func (hs *HashStubs) GetHashesByObject(db *gorm.DB, o *Object) error {
	return db.
		Table("hashes").
		Where(&Hash{ObjectID: &o.ID}).
		Order("created_at").
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
