package domain

import (
	"errors"
)

var (
	ErrTooManyGenerateRequests = errors.New("too many generate requests")
	ErrGenerationTimeout       = errors.New("generation timeout")
	ErrSettingAlreadyExists    = errors.New("setting already exists")
	ErrSettingNotFound         = errors.New("setting not found")
)
