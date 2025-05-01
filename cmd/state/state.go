package state

import (
	"log"

	"github.com/obynonwane/my_blockchain_prototype/cmd/app/handlers/V1/custom"
	"github.com/obynonwane/my_blockchain_prototype/cmd/database"
)

type State struct {
	Model *database.Models
}

func New(models *database.Models) *State {
	return &State{Model: models}
}

func (s *State) CreateUser(data *custom.User) error {
	_, err := s.Model.User.Create(data)
	if err != nil {
		log.Println("error creating user in state:", err)
	}
	return err
}
