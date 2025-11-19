package postgres

import (
	"github.com/eerzho/telegram-ai/pkg/lru"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type Config struct {
	URL           string `env:"POSTGRES_URL,required"`
	STMTCacheSize int    `env:"POSTGRES_STMT_CACHE_SIZE" envDefault:"10"`
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
	return &DB{
		db:        db,
		stmtCache: lru.NewCache[*sqlx.Stmt](cfg.STMTCacheSize),
	}
}
