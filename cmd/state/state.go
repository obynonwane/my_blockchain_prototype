package state

import (
	"log"
	"sync"

	"github.com/obynonwane/my_blockchain_prototype/cmd/app/handlers/V1/custom"
	"github.com/obynonwane/my_blockchain_prototype/cmd/config"
	"github.com/obynonwane/my_blockchain_prototype/cmd/database"
	"github.com/obynonwane/my_blockchain_prototype/cmd/genesis"

	"github.com/obynonwane/my_blockchain_prototype/cmd/mempool"
)

type State struct {
	mu            sync.RWMutex
	beneficiaryID database.AccountID
	evHandler     config.EventHandler
	genesis       genesis.Genesis
	mempool       *mempool.Mempool
	Model         *database.Models
}

// New constructs a new blockchain for data management.
func New(cfg config.Config) (*State, error) {

	// construct mempool with the specified sort strategy
	mempool, err := mempool.NewWithStrategy(cfg.SelectStrategy)
	if err != nil {
		return nil, err
	}

	// Create the State to provide support for managing the blockchain.
	state := &State{
		beneficiaryID: cfg.BeneficiaryID,
		evHandler:     cfg.EvHandler,
		genesis:       cfg.Genesis,
		mempool:       mempool,
		Model:         &cfg.Models,
	}

	return state, nil

}

func (s *State) CreateUser(data *custom.User) error {
	_, err := s.Model.User.Create(data)
	if err != nil {
		log.Println("error creating user in state:", err)
	}

	return err
}

// Shutdown cleanly brings the node down.
func (s *State) Shutdown() error {
	// s.evHandler("state: shutdown: started")
	// defer s.evHandler("state: shutdown: completed")

	return nil
}
