package postgres

import "time"

type Config struct {
	URL             string        `env:"POSTGRES_URL,required"`
	MaxIdleConns    int           `env:"POSTGRES_MAX_IDLE_CONNS"     envDefault:"5"`
	MaxOpenConns    int           `env:"POSTGRES_MAX_OPEN_CONNS"     envDefault:"25"`
	ConnMaxLifeTime time.Duration `env:"POSTGRES_CONN_MAX_LIFE_TIME" envDefault:"10m"`
	ConnMaxIdleTime time.Duration `env:"POSTGRES_CONN_MAX_IDLE_TIME" envDefault:"5m"`
}
