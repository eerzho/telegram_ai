package autootel

import (
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/host"
)

func hostStart() error {
	if err := host.Start(); err != nil {
		return fmt.Errorf("host: %w", err)
	}
	return nil
}
