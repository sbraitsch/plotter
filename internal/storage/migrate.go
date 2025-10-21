package storage

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	// Register the "file" source driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// Register the Postgres database driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

func RunMigrations(databaseURL string) {
	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("✅ Migrations applied successfully")
}
