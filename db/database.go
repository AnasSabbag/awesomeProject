package db

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	"log"
)

func MigrationUp(db *sqlx.DB) {

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	fmt.Println("driver: ", driver)
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations/",
		"postgres", driver)
	if err != nil {
		log.Fatalf("Error creating migration %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Error running migration %v", err)
	}
	log.Println("Migration up completed successfully")
}
