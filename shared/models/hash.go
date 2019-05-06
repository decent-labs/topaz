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

// HashWithApp ...
type HashWithApp struct {
	HashID            string     `gorm:"column:hash_id"`
	HashCreatedAt     time.Time  `gorm:"column:hash_created_at"`
	HashUpdatedAt     time.Time  `gorm:"column:hash_updated_at"`
	HashDeletedAt     *time.Time `gorm:"column:hash_deleted_at"`
	HashMultiHash     []byte     `gorm:"column:hash_multihash"`
	HashUnixTimestamp int64      `gorm:"column:hash_unix_timestamp"`
	HashObjectID      *string    `gorm:"column:hash_object_id"`
	HashProofID       *string    `gorm:"column:hash_proof_id"`
	AppID             string     `gorm:"column:app_id"`
	AppCreatedAt      time.Time  `gorm:"column:app_created_at"`
	AppUpdatedAt      time.Time  `gorm:"column:app_updated_at"`
	AppDeletedAt      *time.Time `gorm:"column:app_deleted_at"`
	AppInterval       int        `gorm:"column:app_interval"`
	AppName           string     `gorm:"column:app_name"`
	AppLastProofed    *int64     `gorm:"column:app_last_proofed"`
	AppUserID         string     `gorm:"column:app_user_id"`
}

// HashesWithApp ...
type HashesWithApp []HashWithApp

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

// GetHashesForProofing ...
func (hs *HashesWithApp) GetHashesForProofing(db *gorm.DB) error {
	clause := `
	SELECT 
		h.id AS hash_id,
		h.created_at AS hash_created_at,
		h.updated_at AS hash_updated_at,
		h.deleted_at AS hash_deleted_at,
		h.hash AS hash_multihash,
		h.unix_timestamp AS hash_unix_timestamp,
		h.object_id AS hash_object_id,
		h.proof_id AS hash_proof_id,
		a.id AS app_id,
		a.created_at AS app_created_at,
		a.updated_at AS app_updated_at,
		a.deleted_at AS app_deleted_at,
		a.interval AS app_interval,
		a.name AS app_name,
		a.last_proofed AS app_last_proofed,
		a.user_id AS app_user_id
	FROM
		hashes h
		JOIN objects o ON o.id = h.object_id
		JOIN apps a ON a.id = o.app_id
	WHERE
		h.proof_id IS NULL
		AND
		CASE
			WHEN a.last_proofed IS NOT NULL THEN (EXTRACT(epoch FROM now()) - a.last_proofed >= a.interval)
			ELSE (EXTRACT(epoch FROM now()) - EXTRACT(epoch FROM a.created_at) >= a.interval)
		END
	`

	return db.Raw(clause).Scan(&hs).Error
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
