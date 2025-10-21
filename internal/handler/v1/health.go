package v1

import "net/http"

type Health struct {
}

func NewHealth() *Health {
	return &Health{}
}

func (h *Health) Check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}
