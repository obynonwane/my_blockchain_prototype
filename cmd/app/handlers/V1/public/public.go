package public

import (
	"net/http"
)

type Handlers struct{}

func (h *Handlers) Genesis(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("public route"))
}
