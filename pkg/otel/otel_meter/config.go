package otelmeter

type Config struct {
	CardinalityLimit int `env:"OTEL_METRIC_CARDINALITY_LIMIT" envDefault:"2000"`
}
