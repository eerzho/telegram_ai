package httphandler

import (
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

func Wrap(f HandlerFunc, e ErrorHandler) http.Handler {
	return http.HandlerFunc(WrapFunc(f, e))
}

func WrapFunc(f HandlerFunc, e ErrorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil && e != nil {
			e(w, r, err)
		}
	}
}
