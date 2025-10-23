package json

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var validate = validator.New()

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("json decode: %w", err)
	}

	if err := validate.Struct(v); err != nil {
		return v, fmt.Errorf("validation: %w", err)
	}

	return v, nil
}

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "json encode error", http.StatusBadRequest)
		panic(err)
	}
}

func EncodeError(w http.ResponseWriter, r *http.Request, status int, err error) {
	Encode(w, r, status, ErrorResponse{Error: err.Error()})
}
