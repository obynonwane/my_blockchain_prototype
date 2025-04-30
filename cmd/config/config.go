package config

import (
	"database/sql"

	"github.com/obynonwane/my_blockchain_prototype/cmd/database"
)

type Config struct {
	DB     *sql.DB
	Models database.Models
}
