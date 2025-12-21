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
	if _, ok := errorhelp.AsType[validator.ValidationErrors](err); ok {
		return slog.LevelInfo
	} else if errorhelp.Any(
		err,
		httpjson.ErrInvalidContentType,
	) {
		return slog.LevelInfo
	} else if errorhelp.Any(
		err,
		domain.ErrGenerationTimeout,
		domain.ErrTooManyGenerateRequests,
	) {
		return slog.LevelWarn
	}
	return slog.LevelError
}

func errorToJSON(err error) httpjson.Error {
	if errors.Is(err, httpjson.ErrInvalidContentType) {
		return httpjson.Error{
			Status:  http.StatusBadRequest,
			Message: http.StatusText(http.StatusBadRequest),
		}
	} else if errors.Is(err, domain.ErrGenerationTimeout) {
		return httpjson.Error{
			Status:  http.StatusRequestTimeout,
			Message: "Please try again later.",
		}
	} else if errors.Is(err, domain.ErrTooManyGenerateRequests) {
		return httpjson.Error{
			Status:  http.StatusTooManyRequests,
			Message: "Please try again later.",
		}
	} else if validationErrs, ok := errorhelp.AsType[validator.ValidationErrors](err); ok {
		return httpjson.Error{
			Status:  http.StatusBadRequest,
			Message: http.StatusText(http.StatusBadRequest),
			Details: validationErrorsToDetails(validationErrs),
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
		return "this field is required"
	case "min":
		return "minimum " + fieldError.Param() + " characters required"
	case "max":
		return "maximum " + fieldError.Param() + " characters allowed"
	case "gt":
		return "must be greater than " + fieldError.Param()
	case "dive":
		return "invalid item in array"
	default:
		return "validation failed"
	}
}
