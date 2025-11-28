package logger

type Config struct {
	AppName    string `env:"APP_NAME,required"`
	AppVersion string `env:"APP_VERSION,required"`
	Level      string `env:"LOGGER_LEVEL"         envDefault:"info"` // debug, info, warn, error
	Format     string `env:"LOGGER_FORMAT"        envDefault:"json"` // text, json
}
