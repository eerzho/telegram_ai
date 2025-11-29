package otelresource

type DetectorType string

const (
	DetectorEnv       DetectorType = "env"
	DetectorHost      DetectorType = "host"
	DetectorContainer DetectorType = "container"
)

type Config struct {
	Detectors []DetectorType `env:"OTEL_RESOURCE_DETECTORS" envDefault:"env" envSeparator:","`
}
