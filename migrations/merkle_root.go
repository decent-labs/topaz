package main

import (
	"fmt"
	"strings"

	"github.com/decentorganization/topaz/shared/database"
	_ "github.com/lib/pq"
	multihash "github.com/multiformats/go-multihash"
	migrate "github.com/rubenv/sql-migrate"
)

// IntermediateProof ...
type IntermediateProof struct {
	ID              string `gorm:"primary_key"`
	MerkleRoot      string
	MerkleRootBytes []byte
}

// TableName ...
func (IntermediateProof) TableName() string {
	return "proofs"
}

func merkleRootMigration(migrations migrate.MigrationSource) (int, error) {
	records, err := migrate.GetMigrationRecords(db, "postgres")
	if err != nil {
		fmt.Println("couldn't get migration records")
		return 0, err
	}

	lastRecord := records[len(records)-1]
	if strings.Compare(lastRecord.Id, "20190502111700_blockchain_explorers.sql") == 0 {
		n1, err := migrate.ExecMax(db, "postgres", migrations, migrate.Up, 1)
		if err != nil {
			fmt.Println("couldn't execute starting merkle root migration")
			return 0, err
		}

		if err := convertRoots(); err != nil {
			return 0, err
		}

		n2, err := migrate.ExecMax(db, "postgres", migrations, migrate.Up, 1)
		if err != nil {
			fmt.Println("couldn't execute finishing merkle root migration")
			return 0, err
		}

		return n1 + n2, nil
	}

	return 0, nil
}

func convertRoots() error {
	var ips []IntermediateProof
	if err := database.Manager.Order("created_at").Find(&ips).Error; err != nil {
		return err
	}

	for _, ip := range ips {
		var mh multihash.Multihash

		mh, err := multihash.FromB58String(ip.MerkleRoot)
		if err != nil {
			return err
		}

		if err := database.Manager.Model(&ip).UpdateColumn("merkle_root_bytes", mh).Error; err != nil {
			return err
		}
	}

	fmt.Println("performed", len(ips), "merkle root conversions")
	return nil
}
