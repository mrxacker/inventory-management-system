package database

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mrxacker/inventory-management-system/services/auth-service/internal/config"
)

func InitDatabase(cfg *config.Config) (*sqlx.DB, error) {
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
