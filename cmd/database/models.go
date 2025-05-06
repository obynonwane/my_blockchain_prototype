package database

import "time"

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserAccount struct {
	ID        	int64     `json:"id"`
	AccountID   string    `json:"account_id"`
	Balance     string    `json:"balance"` // use string to safely handle big numbers (NUMERIC)
	Nonce       int64     `json:"nonce"`
	CodeHash    string    `json:"code_hash"`
	StorageRoot string    `json:"storage_root"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
