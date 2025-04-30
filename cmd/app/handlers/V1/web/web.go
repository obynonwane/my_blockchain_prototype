package web

import (
	"net/http"

	"github.com/obynonwane/my_blockchain_prototype/cmd/database"
)

type Handlers struct {
	Model *database.Models
}

func (h *Handlers) Genesis(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("web route"))
}
