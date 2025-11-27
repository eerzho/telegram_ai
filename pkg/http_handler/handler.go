package httphandler

import (
	"net/http"

	"github.com/eerzho/telegram-ai/pkg/json"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func Wrap(f HandlerFunc) http.Handler {
	return http.HandlerFunc(WrapFunc(f))
}

func WrapFunc(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			json.EncodeError(w, err)
		}
	}
}
