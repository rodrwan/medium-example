package service

import (
	mediumexample "github.com/rodrwan/medium-example"
	"github.com/rodrwan/medium-example/database"
)

// Users represents a User service that interact with database using Database.
type Users struct {
	DBStore *database.Database
}

// GetByID gets user by ID.
func (u *Users) GetByID(id string) (*mediumexample.User, error) {
	user, err := u.DBStore.UsersService.Get(&mediumexample.UserQueryOptions{
		ID: id,
	})
	if err != nil {
		return nil, err
	}

	addr, err := u.DBStore.AddressesService.Get(user.ID)
	if err != nil {
		return nil, err
	}

	user.Address = addr

	return user, nil
}

// Select gets a collection of user.
func (u *Users) Select() ([]*mediumexample.User, error) {
	users, err := u.DBStore.UsersService.Select()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		addr, err := u.DBStore.AddressesService.Get(user.ID)
		if err != nil {
			return nil, err
		}

		user.Address = addr
	}

	return users, nil
}

// Create create a new user.
func (u *Users) Create(user *mediumexample.User) error {
	if err := u.DBStore.UsersService.Create(user); err != nil {
		return err
	}

	user.Address.UserID = user.ID
	if err := u.DBStore.AddressesService.Create(user.Address); err != nil {
		return err
	}

	return nil
}

// NewService expose a new Users service.
func NewService(store *database.Database) mediumexample.UserService {
	return &Users{
		DBStore: store,
	}
}
