package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
)

var (
	ErrInvalidContentType = errors.New("Content-Type must be application/json")
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Decode[T any](r *http.Request) (T, error) {
	var v T

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return v, fmt.Errorf("json decode: %w", ErrInvalidContentType)
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return v, fmt.Errorf("json decode: %w", err)
	}
	
	if mediaType != "application/json" {
		return v, fmt.Errorf("json decode: %w", ErrInvalidContentType)
	}

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("json decode: %w", err)
	}
	return v, nil
}

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func EncodeError(w http.ResponseWriter, r *http.Request, status int, err error) {
	Encode(w, r, status, ErrorResponse{Error: err.Error()})
}
