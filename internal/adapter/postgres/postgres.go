package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type Config struct {
	URL string `env:"POSTGRES_URL,required"`
}

type DB struct {
	db *sqlx.DB
}

func New(cfg Config) *DB {
	db, err := sqlx.Connect("postgres", cfg.URL)
	if err != nil {
		panic(err)
	}
	return &DB{
		db: db,
	}
}
