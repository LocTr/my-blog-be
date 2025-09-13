package migrations

import (
	"context"
	"embed"
	"fmt"
	"log"

	"github.com/LocTr/my-blog-be/database"
	"github.com/uptrace/bun/migrate"
)

var sqlMigrations embed.FS

var Migrations = migrate.NewMigrations()

func init() {
	if err := Migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}
}

func Migrate() {
	db, err := database.DBConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migrator := migrate.NewMigrator(db, Migrations)

	err = migrator.Init(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	group, err := migrator.Migrate(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if group.ID == 0 {
		fmt.Printf("there are no new migrations to run\n")
	} else {
		fmt.Printf("migrated to %s\n", group)
	}
}
