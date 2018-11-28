package models

import "github.com/jinzhu/gorm"

type Object struct {
	gorm.Model

	DataBlob []byte `json:"dataBlob"`
	Hash     string `json:"hash"`

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

func (os *Objects) GetObjectsByHash(db *gorm.DB, o *Object) error {
	return db.Preload("Proof.Batch").Where(&Object{Hash: o.Hash, AppID: o.AppID}).Find(&os).Error
}

func (os Objects) UpdateProof(db *gorm.DB, proofID *uint) error {
	ids := make([]uint, len(os))
	for i, o := range os {
		ids[i] = o.ID
	}
	return db.Model(Object{}).Where("id IN (?)", ids).Updates(Object{ProofID: proofID}).Error
}
