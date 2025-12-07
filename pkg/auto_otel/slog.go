package autootel

import (
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

func NewSlogHandler() slog.Handler {
	return otelslog.NewHandler("")
}
