package public

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/obynonwane/my_blockchain_prototype/cmd/app/handlers/V1/custom"
	"github.com/obynonwane/my_blockchain_prototype/cmd/web"

	"github.com/obynonwane/my_blockchain_prototype/cmd/state"
)

// type Handlers struct {
// 	Model *database.Models
// 	State *state.State
// }

type Handlers struct {
	State *state.State
}
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handlers) Genesis(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	gen := h.State.Genesis()
	web.Respond(ctx, w, gen, http.StatusOK)
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data *custom.User
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}

	err = h.State.CreateUser(data)
	if err != nil {
		log.Println(err, "The error is here")
	}

	// log.Println(res, "created user")
	log.Println(data)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("public route"))
}
