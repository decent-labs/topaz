package migrations

import (
	"fmt"
	"log"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/gobuffalo/packr"
	migrate "github.com/rubenv/sql-migrate"
)

// Attempt ...
func Attempt() {
	migrations := &migrate.PackrMigrationSource{
		Box: packr.NewBox("sql"),
	}

	m, err := merkleRootMigration(migrations)
	if err != nil {
		log.Println("couldn't execute merkle root migration setup scripts")
		log.Fatal(err)
	}

	n, err := migrate.Exec(database.Manager.DB(), "postgres", migrations, migrate.Up)
	if err != nil {
		log.Println("couldn't execute migrations")
		log.Fatal(err)
	}

	n = n + m

	fmt.Printf("Applied %d migrations!\n", n)
}
