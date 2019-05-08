package migrations

import (
	"fmt"
	"strings"

	"github.com/decentorganization/topaz/shared/database"
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
	n := 0

	records, err := migrate.GetMigrationRecords(database.Manager.DB(), "postgres")
	if err != nil {
		return n, err
	}

	lastRecord := records[len(records)-1]
	if strings.Compare(lastRecord.Id, "20190502111700_blockchain_explorers.sql") == 0 {
		n, err = migrate.ExecMax(database.Manager.DB(), "postgres", migrations, migrate.Up, 1)
		if err != nil {
			return n, err
		}
		err = convertRoots()
	}

	return n, err
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
