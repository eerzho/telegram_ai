package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/eerzho/telegram-ai/internal/domain"
	errorhelp "github.com/eerzho/telegram-ai/pkg/error_help"
	httphandler "github.com/eerzho/telegram-ai/pkg/http_handler"
	"github.com/eerzho/telegram-ai/pkg/json"
	"github.com/go-playground/validator/v10"
)

func errorHandler(logger *slog.Logger) httphandler.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		logLevel := errorLogLevel(err)
		logger.Log(r.Context(), logLevel, "handler error", slog.Any("error", err))

		jsonError := errorToJSON(err)
		json.EncodeError(w, jsonError)
	}
}

func errorLogLevel(err error) slog.Level {
	switch {
	case errors.Is(err, json.ErrInvalidContentType):
		return slog.LevelInfo
	case errorhelp.Any(
		err,
		domain.ErrTooManyGenerateRequests,
		domain.ErrGenerationTimeout,
	):
		return slog.LevelWarn
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return slog.LevelInfo
	}

	return slog.LevelError
}

func errorToJSON(err error) json.Error {
	switch {
	case errors.Is(err, json.ErrInvalidContentType):
		return json.Error{
			Status:  http.StatusBadRequest,
			Message: http.StatusText(http.StatusBadRequest),
		}
	case errors.Is(err, domain.ErrTooManyGenerateRequests):
		return json.Error{
			Status:  http.StatusTooManyRequests,
			Message: "Please try again later.",
		}
	case errors.Is(err, domain.ErrGenerationTimeout):
		return json.Error{
			Status:  http.StatusRequestTimeout,
			Message: "Please try again later.",
		}
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return json.Error{
			Status:  http.StatusBadRequest,
			Message: http.StatusText(http.StatusBadRequest),
			Details: validationErrorsToDetails(validationErrors),
		}
	}

	return json.Error{
		Status:  http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}
}

func validationErrorsToDetails(validationErrors validator.ValidationErrors) []json.ErrorDetail {
	details := make([]json.ErrorDetail, 0, len(validationErrors))
	for _, fieldError := range validationErrors {
		message := fieldErrorMessage(fieldError)
		details = append(details, json.ErrorDetail{
			Field:   fieldError.Field(),
			Message: message,
		})
	}
	return details
}

func fieldErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Minimum " + fieldError.Param() + " characters required"
	case "max":
		return "Maximum " + fieldError.Param() + " characters allowed"
	case "gt":
		return "Must be greater than " + fieldError.Param()
	case "dive":
		return "Invalid item in array"
	default:
		return "Validation failed"
	}
}
