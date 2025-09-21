package infrastructure

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"service-template/internal/app/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func MustConfigure(cfg *config.Config) {
	err := runMigrations(cfg)
	if err != nil {
		panic(err)
	}
}

func runMigrations(cfg *config.Config) error {
	db, err := sql.Open("postgres", cfg.Db.Url)
	if err != nil {
		return fmt.Errorf("can not open connection to db: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE SCHEMA IF NOT EXISTS migrations")
	if err != nil {
		return fmt.Errorf("can not create schema for migrations: %w", err)
	}

	dbDriver, err := postgres.WithInstance(
		db, &postgres.Config{
			SchemaName: "migrations",
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	source, err := createEmbedSource()
	if err != nil {
		return fmt.Errorf("can to create embed source: %w", err)
	}

	m, err := migrate.NewWithInstance("embed-source", source, "", dbDriver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	log.Printf("Current migration version: %d (dirty: %t)", version, dirty)

	return nil
}

func createEmbedSource() (source.Driver, error) {
	driver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create iofs source: %w", err)
	}

	return driver, nil
}
