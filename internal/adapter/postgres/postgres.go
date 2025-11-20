package postgres

import (
	"time"

	"github.com/eerzho/telegram-ai/pkg/lru"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type Config struct {
	URL             string        `env:"POSTGRES_URL,required"`
	STMTCacheSize   int           `env:"POSTGRES_STMT_CACHE_SIZE" envDefault:"10"`
	MaxOpenConns    int           `env:"POSTGRES_MAX_OPEN_CONNS"    envDefault:"25"`
	MaxIdleConns    int           `env:"POSTGRES_MAX_IDLE_CONNS"    envDefault:"5"`
	ConnMaxLifetime time.Duration `env:"POSTGRES_CONN_MAX_LIFETIME" envDefault:"5m"`
}

type DB struct {
	db        *sqlx.DB
	stmtCache *lru.Cache[*sqlx.Stmt]
}

func New(cfg Config) *DB {
	db, err := sqlx.Connect("postgres", cfg.URL)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return &DB{
		db:        db,
		stmtCache: lru.NewCache[*sqlx.Stmt](cfg.STMTCacheSize),
	}
}

func (db *DB) Close() error {
	return db.db.Close()
}
