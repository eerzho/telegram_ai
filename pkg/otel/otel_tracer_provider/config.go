package oteltracerprovider

import (
	"time"
)

type Config struct {
	BatchTimeout       time.Duration `env:"OTEL_BSP_SCHEDULE_DELAY" envDefault:"5s"`
	ExportTimeout      time.Duration `env:"OTEL_BSP_EXPORT_TIMEOUT" envDefault:"30s"`
	MaxExportBatchSize int           `env:"OTEL_BSP_MAX_EXPORT_BATCH_SIZE" envDefault:"512"`
	MaxQueueSize       int           `env:"OTEL_BSP_MAX_QUEUE_SIZE" envDefault:"2048"`
}
