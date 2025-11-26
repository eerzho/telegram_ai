package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
)

var (
	ErrInvalidContentType = errors.New("Content-Type must be application/json")
)

type ErrorDetail struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

type Error struct {
	Status  int           `json:"-"`
	Code    string        `json:"code,omitempty"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%#v", e)
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

	if err = json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("json decode: %w", err)
	}
	return v, nil
}

func Encode[T any](w http.ResponseWriter, status int, v T) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
}

func EncodeError(w http.ResponseWriter, err error) {
	var jsonErr Error
	if !errors.As(err, &jsonErr) {
		jsonErr = Error{
			Status:  http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}
	Encode(w, jsonErr.Status, jsonErr)
}
