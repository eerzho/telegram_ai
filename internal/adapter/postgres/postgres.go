package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func New(cfg Config) (*DB, error) {
	db, err := sqlx.Connect("postgres", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifeTime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return &DB{
		DB: db,
	}, nil
}

func MustNew(cfg Config) *DB {
	db, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return db
}
