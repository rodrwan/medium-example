package service

import (
	"context"

	mediumexample "github.com/rodrwan/medium-example"
	"github.com/rodrwan/medium-example/database"
)

// Users represents a User service that interact with database using Database.
type Users struct {
	DBStore *database.Database
}

// GetByID gets user by ID.
func (u *Users) GetByID(ctx context.Context, id string) (*mediumexample.User, error) {
	user, err := u.DBStore.UsersService.Get(ctx, &mediumexample.UserQueryOptions{
		ID: id,
	})
	if err != nil {
		return nil, err
	}

	addr, err := u.DBStore.AddressesService.Get(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	user.Address = addr

	return user, nil
}

// Select gets a collection of user.
func (u *Users) Select(ctx context.Context) ([]*mediumexample.User, error) {
	users, err := u.DBStore.UsersService.Select(ctx)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		addr, err := u.DBStore.AddressesService.Get(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		user.Address = addr
	}

	return users, nil
}

// Create create a new user.
func (u *Users) Create(ctx context.Context, user *mediumexample.User) error {
	if err := u.DBStore.UsersService.Create(ctx, user); err != nil {
		return err
	}

	user.Address.UserID = user.ID
	if err := u.DBStore.AddressesService.Create(ctx, user.Address); err != nil {
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
