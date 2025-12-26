package valkey

import "time"

type Config struct {
	Address  []string      `env:"VALKEY_ADDRESS,required" envSeparator:","`
	Username string        `env:"VALKEY_USERNAME"`
	Password string        `env:"VALKEY_PASSWORD"`
	TTL      time.Duration `env:"VALKEY_TTL"                               envDefault:"1h"`
}
