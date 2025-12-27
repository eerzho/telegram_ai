package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	httpjson "github.com/eerzho/goiler/pkg/http_json"
	httpserver "github.com/eerzho/goiler/pkg/http_server"
	"github.com/eerzho/telegram_ai/internal/domain"
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
	}
	if _, ok := errorhelp.AsType[*json.UnmarshalTypeError](err); ok {
		return slog.LevelInfo
	}
	if _, ok := errorhelp.AsType[*json.SyntaxError](err); ok {
		return slog.LevelInfo
	}

	if errorhelp.Any(
		err,
		httpjson.ErrInvalidContentType,
		domain.ErrSettingAlreadyExists,
		domain.ErrSettingNotFound,
	) {
		return slog.LevelInfo
	}

	if errorhelp.Any(
		err,
		domain.ErrGenerationTimeout,
		domain.ErrTooManyGenerateRequests,
	) {
		return slog.LevelWarn
	}

	return slog.LevelError
}

func errorToJSON(err error) httpjson.Error {
	if validationErr := parseValidationError(err); validationErr != nil {
		return *validationErr
	}

	if unmarshalErr := parseUnmarshalTypeError(err); unmarshalErr != nil {
		return *unmarshalErr
	}

	if syntaxErr := parseSyntaxError(err); syntaxErr != nil {
		return *syntaxErr
	}

	if knownErr := checkKnownError(err); knownErr != nil {
		return *knownErr
	}

	return httpjson.Error{
		Status:  http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}
}

func parseValidationError(err error) *httpjson.Error {
	validationErrs, ok := errorhelp.AsType[validator.ValidationErrors](err)
	if !ok {
		return nil
	}

	return &httpjson.Error{
		Status:  http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
		Details: validationErrorsToDetails(validationErrs),
	}
}

func parseUnmarshalTypeError(err error) *httpjson.Error {
	unmarshalTypeErr, ok := errorhelp.AsType[*json.UnmarshalTypeError](err)
	if !ok {
		return nil
	}

	message := http.StatusText(http.StatusBadRequest)
	if unmarshalTypeErr.Field != "" {
		message = "invalid type for field " + unmarshalTypeErr.Field
	}

	return &httpjson.Error{
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

func parseSyntaxError(err error) *httpjson.Error {
	_, ok := errorhelp.AsType[*json.SyntaxError](err)
	if !ok {
		return nil
	}

	return &httpjson.Error{
		Status:  http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
	}
}

func checkKnownError(err error) *httpjson.Error {
	mappings := map[error]httpjson.Error{
		httpjson.ErrInvalidContentType: {
			Status:  http.StatusBadRequest,
			Message: httpjson.ErrInvalidContentType.Error(),
		},
		domain.ErrGenerationTimeout: {
			Status:  http.StatusRequestTimeout,
			Message: domain.ErrGenerationTimeout.Error(),
		},
		domain.ErrTooManyGenerateRequests: {
			Status:  http.StatusTooManyRequests,
			Message: domain.ErrTooManyGenerateRequests.Error(),
		},
		domain.ErrSettingAlreadyExists: {
			Status:  http.StatusConflict,
			Message: domain.ErrSettingAlreadyExists.Error(),
		},
		domain.ErrSettingNotFound: {
			Status:  http.StatusNotFound,
			Message: domain.ErrSettingNotFound.Error(),
		},
	}

	for knownErr, mapping := range mappings {
		if knownErr != nil && errors.Is(err, knownErr) {
			return &mapping
		}
	}

	return nil
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
