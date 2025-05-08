package state

import (
	"log"
	"strconv"
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

// Genesis returns a copy of the genesis information
func (s *State) Genesis() genesis.Genesis {
	return s.genesis
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

// MempoolLength returns the current length of the mempool
func (s *State) MempoolLength() int {
	return s.mempool.Count()
}

// Mempool returns a copy of the mempool
func (s *State) Mempool() []database.BlockTx {
	return s.mempool.PickBest()
}

// UpsertMempool adds a new transaction to the mempool
func (s *State) UpsertMempool(tx database.BlockTx) error {
	return s.mempool.Upsert(tx)
}

func (s *State) Accounts() (map[database.AccountID]database.Account, error) {
	data, err := s.Model.UserAccount.Copy()
	if err != nil {
		s.evHandler("state: error retrieving accounts copy", err)
	}

	// convert the data to a map
	accounts := make(map[database.AccountID]database.Account)

	for _, ua := range data {
		balance, err := strconv.ParseUint(ua.Balance, 10, 64)
		if err != nil {
			s.evHandler("state: error converting balance to parseUint", err)
		}

		account := database.Account{
			AccountID: database.AccountID(ua.AccountID),
			Nonce:     uint64(ua.Nonce),
			Balance:   balance,
		}

		accounts[database.AccountID(ua.AccountID)] = account
	}
	return accounts, nil
}

func (s *State) SeedGenesisAccount(genesis *genesis.Genesis) error {

	err := s.Model.UserAccount.SeedGenesisAccount(genesis)
	if err != nil {
		s.evHandler("state: error seeding genesis accounts", err)
	}

	return err
}
