package errorhelp

import (
	"errors"
	"fmt"
)

func Any(err error, targets ...error) bool {
	for _, target := range targets {
		if is := errors.Is(err, target); is {
			return true
		}
	}
	return false
}

func WithOP(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}
