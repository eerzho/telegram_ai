package domain

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/eerzho/telegram-ai/pkg/json"
	"github.com/go-playground/validator/v10"
)

var (
	ErrSummaryNotFound         = errors.New("summary not found")
	ErrTooManyGenerateRequests = errors.New("too many generate requests")
	ErrGenerationTimeout       = errors.New("generation timeout")
)

func LogLevel(err error) slog.Level {
	switch {
	case errors.Is(err, ErrSummaryNotFound):
		return slog.LevelInfo
	case errors.Is(err, ErrTooManyGenerateRequests):
		return slog.LevelWarn
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return slog.LevelInfo
	}

	return slog.LevelError
}

func MapToJSONError(err error) json.Error {
	switch {
	case errors.Is(err, ErrSummaryNotFound):
		return json.Error{
			Status:  http.StatusNotFound,
			Message: "Summary not found",
		}
	case errors.Is(err, ErrTooManyGenerateRequests):
		return json.Error{
			Status:  http.StatusTooManyRequests,
			Message: "Please try again later",
		}
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return json.Error{
			Status:  http.StatusBadRequest,
			Message: http.StatusText(http.StatusBadRequest),
			Details: mapToJSONErrorDetails(validationErrors),
		}
	}

	return json.Error{
		Status:  http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}
}

func mapToJSONErrorDetails(validationErrors validator.ValidationErrors) []json.ErrorDetail {
	details := make([]json.ErrorDetail, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		message := validationErrorMessage(validationError)
		details = append(details, json.ErrorDetail{
			Field:   validationError.Field(),
			Message: message,
		})
	}
	return details
}

func validationErrorMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Minimum " + fieldErr.Param() + " characters required"
	case "max":
		return "Maximum " + fieldErr.Param() + " characters allowed"
	case "gt":
		return "Must be greater than " + fieldErr.Param()
	case "dive":
		return "Invalid item in array"
	default:
		return "Validation failed"
	}
}
