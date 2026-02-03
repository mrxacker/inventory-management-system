package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mrxacker/inventory-management-system/services/auth-service/internal/config"
)

func InitDatabase(cfg *config.Config) (*sqlx.DB, error) {
	if err := createDatabaseIfNotExists(cfg); err != nil {
		return nil, err
	}

	db, err := sqlx.Connect("postgres", cfg.Postgres.DBUrl)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	db.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Postgres.ConnMaxLifetime) * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func createDatabaseIfNotExists(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer db.Close()

	// Check if database exists
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = db.QueryRow(query, cfg.Postgres.Name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	if !exists {
		// Create database
		createDBQuery := fmt.Sprintf("CREATE DATABASE %s", cfg.Postgres.Name)
		_, err = db.Exec(createDBQuery)
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}
	return nil
}
