package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/obynonwane/my_blockchain_prototype/cmd/app/handlers/V1/custom"
)

// db timeout period
const dbTimeout = time.Second * 3

// data of sqlDB type here connections to DB will live
var db *sql.DB

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	User User
}

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
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
