package main

import (
	"fmt"
	"log"
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

func merkleRootMigration(migrations migrate.MigrationSource) {
	records, err := migrate.GetMigrationRecords(db, "postgres")
	if err != nil {
		log.Println("couldn't get migration records")
	}

	lastRecord := records[len(records)-1]
	if strings.Compare(lastRecord.Id, "20190502111700_blockchain_explorers.sql") == 0 {
		n, err := migrate.ExecMax(db, "postgres", migrations, migrate.Up, 1)
		if err != nil {
			log.Println("couldn't execute starting merkle root migration")
			log.Fatal(err)
		}

		fmt.Println("executed first", n, "merkle root migration")

		convertRoots()

		n, err = migrate.ExecMax(db, "postgres", migrations, migrate.Up, 1)
		if err != nil {
			log.Println("couldn't execute finishing merkle root migration")
			log.Fatal(err)
		}

		fmt.Println("executed final", n, "merkle root migration")
	}
}

func convertRoots() {
	var intermediateProofs []IntermediateProof
	database.Manager.Order("created_at").Find(&intermediateProofs)

	for _, ip := range intermediateProofs {
		var mh multihash.Multihash

		mh, err := multihash.FromB58String(ip.MerkleRoot)
		if err != nil {
			fmt.Println("oops decode didn't work")
		}

		database.Manager.Model(&ip).UpdateColumn("merkle_root_bytes", mh)
	}
}
