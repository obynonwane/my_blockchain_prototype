package config

import (
	"database/sql"
	"time"

	"github.com/obynonwane/my_blockchain_prototype/cmd/database"
	"github.com/obynonwane/my_blockchain_prototype/cmd/genesis"
)

// EventHandler defines a function that is called when events
// occur in the processing of persisting blocks.
type EventHandler func(v string, args ...any) // represent a function that we can call if we want to log

// config represent the configuration required
// to start the node or bl
type Config struct {
	DB              *sql.DB
	Models          database.Models
	BeneficiaryID   database.AccountID // receiver of mining reward/gas fee for this node
	Genesis         genesis.Genesis
	SelectStrategy  string
	EvHandler       EventHandler  // logging
	ReadTimeout     time.Duration // The maximum time the server will wait to read the entire request from the client.
	WriteTimeout    time.Duration // The maximum time the server will wait to write a response to the client.
	IdleTimeout     time.Duration // maximum time the server will keep a connection open when it’s idle
	ShutdownTimeout time.Duration // The maximum time the server will wait for ongoing requests to finish when it is shutting down.
}
