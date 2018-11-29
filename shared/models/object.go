package models

import (
	"bytes"
	"crypto/sha256"

	"github.com/cbergoon/merkletree"
	"github.com/jinzhu/gorm"

	multihash "github.com/multiformats/go-multihash"
)

type Object struct {
	gorm.Model

	DataBlob      []byte `json:"dataBlob"`
	Hash          string `json:"hash"`
	UnixTimestamp int64  `json:"unixTimestamp"`

	AppID   uint   `json:"appId"`
	App     *App   `json:"app,omitempty"`
	ProofID *uint  `json:"proofId"`
	Proof   *Proof `json:"proof,omitempty"`
}

type Objects []Object

func (o Object) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(o.DataBlob); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func (o Object) Equals(other merkletree.Content) (bool, error) {
	return bytes.Compare(o.DataBlob, other.(Object).DataBlob) == 0, nil
}

func (o Object) MakeHash() (string, error) {
	digest, err := o.CalculateHash()
	if err != nil {
		return "", err
	}

	return getReadableHash(digest)
}

func (os Objects) GetMerkleRoot() (string, error) {
	t, err := makeMerkleTree(&os)
	if err != nil {
		return "", err
	}

	return getReadableHash(t.MerkleRoot())
}

func (o *Object) CreateObject(db *gorm.DB) error {
	return db.Create(&o).Error
}

func (os *Objects) GetObjectsByAppID(db *gorm.DB, id uint) error {
	clause := "proof_id IS NULL AND app_id = ?"
	return db.Where(clause, id).Find(&os).Error
}

func (os *Objects) GetVerifiedObjects(db *gorm.DB, o *Object) error {
	return db.Preload("Proof.Batch").Where(o).Find(&os).Error
}

func (os Objects) UpdateProof(db *gorm.DB, proofID *uint) error {
	ids := make([]uint, len(os))
	for i, o := range os {
		ids[i] = o.ID
	}
	return db.Model(Object{}).Where("id IN (?)", ids).Updates(Object{ProofID: proofID}).Error
}

func (os *Objects) GetObjectsByTimestamps(db *gorm.DB, appId uint, start int, end int) error {
	return db.
		Preload("Proof.Batch").
		Where("app_id = ?", appId).
		Where("unix_timestamp BETWEEN (?) AND (?)", start, end).
		Find(&os).
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

func makeMerkleTree(os *Objects) (*merkletree.MerkleTree, error) {
	var list []merkletree.Content
	for _, obj := range *os {
		list = append(list, obj)
	}

	return merkletree.NewTree(list)
}
