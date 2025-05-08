package state

import (
	"log"

	"github.com/obynonwane/my_blockchain_prototype/cmd/database"
)

func (s *State) QueryAccount(account database.AccountID) (database.UserAccount, error) {
	data, err := s.Model.UserAccount.Query(account)
	if err != nil {
		log.Println("error creating user in state:", err)
	}

	return data, err
}
