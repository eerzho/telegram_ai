package http

import (
	"errors"
	"log/slog"
	"net/http"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	httpjson "github.com/eerzho/goiler/pkg/http_json"
	httpserver "github.com/eerzho/goiler/pkg/http_server"
	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/go-playground/validator/v10"
)

func errorHandler(logger *slog.Logger) httpserver.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		logLevel := errorLogLevel(err)
		logger.Log(r.Context(), logLevel, "request failed", slog.Any("error", err))

		jsonError := errorToJSON(err)
		httpjson.EncodeError(w, jsonError)
	}
}

func errorLogLevel(err error) slog.Level {
	switch {
	case errors.Is(err, httpjson.ErrInvalidContentType):
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

func errorToJSON(err error) httpjson.Error {
	switch {
	case errors.Is(err, httpjson.ErrInvalidContentType):
		return httpjson.Error{
			Status:  http.StatusBadRequest,
			Message: http.StatusText(http.StatusBadRequest),
		}
	case errors.Is(err, domain.ErrTooManyGenerateRequests):
		return httpjson.Error{
			Status:  http.StatusTooManyRequests,
			Message: "Please try again later.",
		}
	case errors.Is(err, domain.ErrGenerationTimeout):
		return httpjson.Error{
			Status:  http.StatusRequestTimeout,
			Message: "Please try again later.",
		}
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return httpjson.Error{
			Status:  http.StatusBadRequest,
			Message: http.StatusText(http.StatusBadRequest),
			Details: validationErrorsToDetails(validationErrors),
		}
	}

	return httpjson.Error{
		Status:  http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}
}

func validationErrorsToDetails(validationErrors validator.ValidationErrors) []httpjson.ErrorDetail {
	details := make([]httpjson.ErrorDetail, 0, len(validationErrors))
	for _, fieldError := range validationErrors {
		message := fieldErrorMessage(fieldError)
		details = append(details, httpjson.ErrorDetail{
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
