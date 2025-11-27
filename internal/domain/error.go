package domain

import (
	"errors"
)

var (
	ErrTooManyGenerateRequests = errors.New("too many generate requests")
	ErrGenerationTimeout       = errors.New("generation timeout")
)
