package public

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/obynonwane/my_blockchain_prototype/cmd/database"
)

type Handlers struct {
	Model database.Models
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handlers) Genesis(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("public route"))
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data *database.User
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}

	err = h.Model.User.Create(data)
	// if err != nil {
	// 	log.Println(err, "The error is here")
	// }

	log.Println("created user")
	log.Println(data)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("public route"))
}
