package database

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

// Account represents information stored in the database for an individual account.
type Account struct {
	AccountID AccountID
	Nonce     uint64
	Balance   uint64
}

// newAccount constructs a new account value for use.
func newAccount(accountID AccountID, balance uint64) Account {
	return Account{
		AccountID: accountID,
		Balance:   balance,
	}
}

// =============================================================================

// AccountID represents an account id that is used to sign transactions and is
// associated with transactions on the blockchain. This will be the last 20
// bytes of the public key.
type AccountID string

// PublicKeyToAccountID converts the public key to an account value.
func PublicKeyToAccountID(pk ecdsa.PublicKey) AccountID {
	return AccountID(crypto.PubkeyToAddress(pk).String())
}
