package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/obynonwane/my_blockchain_prototype/cmd/app/handlers/V1/custom"
	"github.com/obynonwane/my_blockchain_prototype/cmd/genesis"
)

// db timeout period
const dbTimeout = time.Second * 3

// data of sqlDB type here connections to DB will live
var db *sql.DB

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	User        User
	UserAccount UserAccount
}

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User:        User{},
		UserAccount: UserAccount{},
	}
}

func (u *User) Create(data *custom.User) (User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var i User
	stmt := `insert into users (name, email)
		values ($1, $2) RETURNING id, name, email, updated_at, created_at`

	err := db.QueryRowContext(ctx, stmt,
		data.Name,
		data.Email,
	).Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.UpdatedAt,
		&i.CreatedAt,
	)

	if err != nil {
		log.Println(err)
	}
	return i, nil
}

func (ua *UserAccount) SeedGenesisAccount(genesis *genesis.Genesis) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	for account_id, balance := range genesis.Balances {
		log.Printf("Seeding Address: %s with Balance: %v\n", account_id, balance)

		// Check if account exists
		query := `SELECT account_id FROM accounts WHERE account_id = $1`
		var existingID string
		err := db.QueryRowContext(ctx, query, account_id).Scan(&existingID)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// Account doesn't exist, insert it
				insert := `INSERT INTO accounts (account_id, balance, nonce, code_hash, storage_root, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`

				_, err := db.ExecContext(ctx, insert, account_id, balance, 0, "", "")
				if err != nil {
					return fmt.Errorf("failed to insert genesis account %s: %w", account_id, err)
				}
				log.Printf("Inserted genesis account %s\n", account_id)
			} else {
				return fmt.Errorf("error checking existing account %s: %w", account_id, err)
			}
		} else {
			log.Printf("Account %s already exists, skipping\n", account_id)
		}
	}

	log.Println("Genesis seeding completed.")
	return nil
}
