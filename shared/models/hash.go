package models

import (
	"bytes"
	"encoding/hex"

	"github.com/cbergoon/merkletree"
	"github.com/jinzhu/gorm"

	multihash "github.com/multiformats/go-multihash"
)

type Hash struct {
	gorm.Model

	HashHex       string `json:"hash" gorm:"-"`
	Hash          []byte `json:"-"`
	UnixTimestamp int64  `json:"unixTimestamp"`

	ObjectID *uint   `json:"objectId"`
	Object   *Object `json:"object,omitempty"`
	ProofID  *uint   `json:"proofId"`
	Proof    *Proof  `json:"proof,omitempty"`
}

type Hashes []Hash

func (h Hash) CalculateHash() ([]byte, error) {
	return h.Hash, nil
}

func (h Hash) Equals(other merkletree.Content) (bool, error) {
	return bytes.Compare(h.Hash, other.(Hash).Hash) == 0, nil

}

func (hs Hashes) GetMerkleRoot() (string, error) {
	t, err := makeMerkleTree(&hs)
	if err != nil {
		return "", err
	}

	return getReadableHash(t.MerkleRoot())
}

func (h *Hash) CreateHash(db *gorm.DB) error {
	return db.Create(&h).Error
}

func (hs *Hashes) GetHashesByApp(db *gorm.DB, app *App) error {
	return db.Model(&app).
		Table("hashes").
		Joins("join objects on objects.id = hashes.object_id").
		Where("hashes.proof_id IS NULL").
		Find(&hs).
		Error
}

func (hs *Hashes) GetVerifiedHashes(db *gorm.DB, a *App, h *Hash) error {
	h.Hash, _ = hex.DecodeString(h.HashHex)

	if err := db.Model(&a).
		Table("hashes").
		Joins("join objects on objects.id = hashes.object_id").
		Where(h).
		Find(&hs).Error; err != nil {
		return err
	}

	for _, h := range *hs {
		if h.Proof == nil {
			continue
		}

		if err := h.Proof.CheckValidity(); err != nil {
			return err
		}
	}

	return nil
}

func (hs Hashes) UpdateProof(db *gorm.DB, proofID *uint) error {
	ids := make([]uint, len(hs))
	for i, h := range hs {
		ids[i] = h.ID
	}

	return db.Model(Hash{}).Where("id IN (?)", ids).Updates(Hash{ProofID: proofID}).Error
}

func (hs *Hashes) GetHashesByTimestamps(db *gorm.DB, appID uint, start int, end int) error {
	return db.
		Preload("Proof.Batch").
		Where("app_id = ?", appID).
		Where("unix_timestamp BETWEEN (?) AND (?)", start, end).
		Find(&hs).
		Error
}

func getReadableHash(digest []byte) (string, error) {
	mhBuf, err := multihash.Encode(digest, multihash.SHA2_256)
	if err != nil {
		return "", err
	}

	mh, err := multihash.Cast(mhBuf)
	if err != nil {
		return "", err
	}

	return mh.B58String(), nil
}

func makeMerkleTree(hs *Hashes) (*merkletree.MerkleTree, error) {
	var list []merkletree.Content
	for _, obj := range *hs {
		list = append(list, obj)
	}

	return merkletree.NewTree(list)
}
