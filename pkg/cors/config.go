package cors

type Config struct {
	AllowedOrigins   string `env:"CORS_ALLOWED_ORIGINS"   envDefault:"*"`
	AllowedMethods   string `env:"CORS_ALLOWED_METHODS"   envDefault:"GET,POST,PUT,DELETE,OPTIONS,PATCH"`
	AllowedHeaders   string `env:"CORS_ALLOWED_HEADERS"   envDefault:"Content-Type,Authorization,X-Request-Id"`
	AllowCredentials bool   `env:"CORS_ALLOW_CREDENTIALS" envDefault:"false"`
	MaxAge           int    `env:"CORS_MAX_AGE"           envDefault:"3600"`
}
