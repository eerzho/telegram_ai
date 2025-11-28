package bodysize

type Config struct {
	Max int `env:"BODY_SIZE_MAX" envDefault:"5"`
}
