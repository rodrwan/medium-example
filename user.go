package mediumexample

import (
	"errors"
	"time"
)

// User hold user information returned by UserService
type User struct {
	ID    string `json:"id,omitempty" db:"id"`
	Email string `json:"email,omitempty" db:"email"`

	FirstName string    `json:"first_name,omitempty" db:"first_name"`
	LastName  string    `json:"last_name,omitempty" db:"last_name"`
	Phone     string    `json:"phone,omitempty" db:"phone"`
	Birthdate time.Time `json:"birthdate,omitempty" db:"birthdate"`

	Address *Address `json:"address,omitempty" db:"-"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Users error definitions.
var (
	ErrUserNotFound       = errors.New("user: user not found")
	ErrUserAlreadyExists  = errors.New("user: user already exists")
	ErrUserInvalidAddress = errors.New("user: invalid address fields")
	ErrAddressNotFound    = errors.New("user: address not found")
)

// UserQueryOptions represents query options to filter users.
type UserQueryOptions struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
}

// UserService define service behavior to operate Users.
// In the implementation of this interface we use Database
// that let us operate over users and addresses service.
type UserService interface {
	GetByID(string) (*User, error)
	Select() ([]*User, error)

	Create(*User) error
}
