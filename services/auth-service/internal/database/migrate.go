package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/mrxacker/inventory-management-system/services/auth-service/internal/config"
)

func RunMigrations(cfg *config.Config) error {

	db, err := sql.Open("postgres", cfg.Postgres.DBUrl)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get version: %w", err)
	}

	log.Printf("Current migration version: %d (dirty: %v)", version, dirty)

	if dirty {
		return fmt.Errorf(
			"database is in dirty migration state at version %d; "+
				"manual intervention required (fix migration and run migrate force)",
			version,
		)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	version, _, _ = m.Version()
	log.Printf("Migration completed. Current version: %d", version)

	return nil
}
