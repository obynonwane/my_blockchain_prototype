package database

import (
	"database/sql"
	"time"
)

// db timeout period
const dbTimeout = time.Second * 3

// data of sqlDB type here connections to DB will live
var db *sql.DB

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	User User
}

type User struct {
	ID          int64     `json:"id"`
	FirstName   string    `json:"first_name,omitempty"`
	LastName    string    `json:"last_name,omitempty"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserType    string    `json:"user_type"`
	UserSubType string    `json:"user_sub_type"`
	IsApproved  bool      `json:"is_approved"`
	IsActive    bool      `json:"is_active"`
}
