package database

import (
	"database/sql"

	"github.com/bobbybaiOuO/BShortUrl/config"
	_ "github.com/lib/pq"
)

// NewDB .
func NewDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}