package mempool

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/obynonwane/my_blockchain_prototype/cmd/database"
	"github.com/obynonwane/my_blockchain_prototype/cmd/mempool/selector"
)

// Mempool represents a cache of transactions organised by account:nonce
type Mempool struct {
	mu       sync.RWMutex
	pool     map[string]database.BlockTx // string here which is the key is the concate of from addresss and nonce
	selectFn selector.Func
}

// New constructs a new mempool using the default sort strategy
func New() (*Mempool, error) {
	return NewWithStrategy(selector.StrategyTip)
}

// NewWithStrategy constructs a new mempool with specified sort strategy
func NewWithStrategy(strategy string) (*Mempool, error) {
	selectFn, err := selector.Retrieve(strategy)
	if err != nil {
		return nil, err
	}
	mp := Mempool{
		pool:     make(map[string]database.BlockTx),
		selectFn: selectFn,
	}
	return &mp, nil
}

// Count returns the current number of transactions in the pool.
func (mp *Mempool) Count() int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return len(mp.pool)
}

// Upsert adds or replace a transaction from th mempool
func (mp *Mempool) Upsert(tx database.BlockTx) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	// CORE NOTE: Different blockchains have different algorithms to limit the
	// size of the mempool. Some limit based on the amount of memory being
	// consumed and some may limit based on the number of transaction. If a limit
	// is met, then either the transaction that has the least return on investment
	// or the oldest will be dropped from the pool to make room for new the transaction.

	// For now my blockchain is not imposing any limits
	key, err := mapKey(tx)
	if err != nil {
		return err
	}

	// Ethereum requires a 10% bump in the tip to replace an existing
	// transaction in the mempool and so do we. We want to limit users
	// from this sort of behavior.
	if etx, exists := mp.pool[key]; exists {
		if tx.Tip < uint64(math.Round(float64(etx.Tip)*1.10)) {
			return errors.New("replacing a transaction requires a 10% bump in the tip")
		}
	}

	// update the tx in the mempool
	// the key is the concate with (:) of from address and nonce
	mp.pool[key] = tx
	return nil
}

// Delete removed a transaction from the mempool
func (mp *Mempool) Delete(tx database.BlockTx) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	key, err := mapKey(tx)
	if err != nil {
		return err
	}

	// delete item from the pool
	delete(mp.pool, key)

	return nil
}

// PickBest uses the configured sort strategy to return a set of transactions.
// If 0 is passed, all transactions in the mempool will be returned.
func (mp *Mempool) PickBest(howMany ...uint16) []database.BlockTx {
	number := 0
	if len(howMany) > 0 {
		number = int(howMany[0])
	}

	// CORE NOTE: Most blockchains do set a max block size limit and this size
	// will determined which transactions are selected. When picking the best
	// transactions for the next block, the Ardan blockchain is currently not
	// focused on block size but a max number of transactions.
	//
	// When the selection algorithm does need to consider sizing, picking the
	// right transactions that maximize profit gets really hard. On top of this,
	// today a miner gets a mining reward for each mined block. In the future
	// this could go away leaving just fees for the transactions that are
	// selected as the only form of revenue. This will change how transactions
	// need to be selected.

	// copy all the transactions for each account into a separate slices.
	m := make(map[database.AccountID][]database.BlockTx) // map where value are slices
	mp.mu.RLock()
	{
		if number == 0 {
			number = len(mp.pool)
		}
		// loop through the mempool and select transaction
		// and append it to m slice
		for key, tx := range mp.pool {
			account := accountFromMapKey(key)
			m[account] = append(m[account], tx)
		}
	}
	mp.mu.RUnlock()

	//The selection algorithms is expecting this slice of transactions
	// organised by account
	return mp.selectFn(m, number)

}

// Truncate clears all the transaction from the pool
func (mp *Mempool) Truncate() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.pool = make(map[string]database.BlockTx)
}

// ========================================================================================

// mapKey is used to generate the map key.
func mapKey(tx database.BlockTx) (string, error) {
	return fmt.Sprintf("%s:%d", tx.FromID, tx.Nonce), nil
}

// accountFromMapKey extracts the account information from the mapKey
func accountFromMapKey(key string) database.AccountID {
	// splits key string and convert to accountID type returning th 0 indexed
	return database.AccountID(strings.Split(key, ":")[0])
}
