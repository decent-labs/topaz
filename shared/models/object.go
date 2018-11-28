package models

import (
	"github.com/jinzhu/gorm"
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

func (o *Object) CreateObject(db *gorm.DB) error {
	return db.Create(&o).Error
}

func (os *Objects) GetObjectsByAppID(db *gorm.DB, id uint) error {
	clause := "proof_id IS NULL AND app_id = ?"
	return db.Where(clause, id).Find(&os).Error
}

func (os *Objects) GetObjects(db *gorm.DB, o *Object) error {
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
