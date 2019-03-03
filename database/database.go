package database

import (
	"os"

	"github.com/jmoiron/sqlx"
	mediumexample "github.com/rodrwan/medium-example"
	"github.com/rodrwan/medium-example/database/postgres"
)

// UsersService is the interface that represent basic operation over
// users table in any database.
type UsersService interface {
	Get(*mediumexample.UserQueryOptions) (*mediumexample.User, error)
	Select() ([]*mediumexample.User, error)

	Create(*mediumexample.User) error
	Update(*mediumexample.User) error
	Delete(*mediumexample.User) error
}

// AddressesService is the interface that represent basic operation over
// addresses table in any database.
type AddressesService interface {
	Get(string) (*mediumexample.Address, error)

	Create(*mediumexample.Address) error
	Update(*mediumexample.Address) error
	Delete(*mediumexample.Address) error
}

// Database hold the database services.
type Database struct {
	UsersService     UsersService
	AddressesService AddressesService
}

// NewPostgres creates a new Database with postgres as driver.
func NewPostgres(dsn string) (*Database, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &Database{
		UsersService: &postgres.UsersService{
			Store:  db,
			Logger: postgres.NewDBLogger(os.Stdout, "users"),
		},
		AddressesService: &postgres.AddressesService{
			Store:  db,
			Logger: postgres.NewDBLogger(os.Stdout, "addresses"),
		},
	}, nil
}
