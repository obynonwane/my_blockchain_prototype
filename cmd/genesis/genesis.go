package genesis

import (
	"encoding/json"
	"os"
	"time"
)

// Genesis represents the genesis file.
type Genesis struct {
	Date          time.Time         `json:"date"`
	ChainID       uint16            `json:"chain_id"`        // The chain id represents an unique id for this running instance.
	TransPerBlock uint16            `json:"trans_per_block"` // The maximum number of transactions that can be in a block.
	Difficulty    uint16            `json:"difficulty"`      // How difficult it needs to be to solve the work problem.
	MiningReward  uint64            `json:"mining_reward"`   // Reward for mining a block.
	GasPrice      uint64            `json:"gas_price"`       // Fee paid for each transaction mined into a block.
	Balances      map[string]uint64 `json:"balances"`
}

// function to load Genesis.json file
func Load() (Genesis, error) {
	// specify path to the genesis file
	path := "cmd/zblock/genesis.json"

	// reads the json given into a slice of byte
	content, err := os.ReadFile(path)
	if err != nil {
		return Genesis{}, err
	}

	var genesis Genesis

	// converting back the slice byte above
	// into go struct
	err = json.Unmarshal(content, &genesis)
	if err != nil {
		return Genesis{}, err
	}

	return genesis, nil
}
