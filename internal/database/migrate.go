package database

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) {
	// Try common paths for migrations
	path := "file://migrations"
	m, err := migrate.New(path, dbURL)
	if err != nil {
		// Try fallback to interaction_service/migrations
		path = "file://interaction_service/migrations"
		m, err = migrate.New(path, dbURL)
		if err != nil {
			log.Printf("Could not create migrate instance: %v", err)
			return
		}
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Could not run up migrations: %v", err)
		return
	}

	log.Println("Migrations applied successfully!")
}
