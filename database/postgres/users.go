package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	mediumexample "github.com/rodrwan/medium-example"
	"github.com/sirupsen/logrus"
)

var (
	// ErrUserNotFound user not found.
	ErrUserNotFound = errors.New("user not found")
)

// UsersService is a service that query users table using SQLExecutor.
type UsersService struct {
	Store  SQLExecutor
	Logger *logrus.Logger
}

// Get gets a user by query params.
func (us *UsersService) Get(query *mediumexample.UserQueryOptions) (*mediumexample.User, error) {
	q := squirrel.Select("*").From("users").Where("deleted_at is null")

	if query.ID != "" {
		q = q.Where("id = ?", query.ID)
	}

	sqlString, args, err := q.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var u mediumexample.User
	row := us.Store.QueryRowx(sqlString, args...)
	if err := row.StructScan(&u); err != nil {
		return nil, us.userError(sqlString, args, err)
	}

	us.Logger.WithFields(logrus.Fields{
		"query": sqlString,
	}).Info("OK")

	return &u, nil
}

// Select gets a collection of user.
func (us *UsersService) Select() ([]*mediumexample.User, error) {
	q := squirrel.Select("*").From("users").Where("deleted_at is null")
	sqlString, args, err := q.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := us.Store.Queryx(sqlString, args...)
	if err != nil {
		return nil, err
	}

	uu := make([]*mediumexample.User, 0)

	for rows.Next() {
		var u mediumexample.User
		if err := rows.StructScan(&u); err != nil {
			return nil, us.userError(sqlString, args, err)
		}

		uu = append(uu, &u)
	}

	us.Logger.WithFields(logrus.Fields{
		"query": sqlString,
	}).Info("OK")

	return uu, nil
}

// Create creates a new user.
func (us *UsersService) Create(u *mediumexample.User) error {
	sqlString, args, err := squirrel.Insert("users").
		Columns("email", "first_name", "last_name", "phone", "birthdate").
		Values(u.Email, u.FirstName, u.LastName, u.Phone, u.Birthdate).
		Suffix("returning *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	row := us.Store.QueryRowx(sqlString, args...)
	if err := row.StructScan(u); err != nil {
		return us.userError(sqlString, args, err)
	}

	us.Logger.WithFields(logrus.Fields{
		"query": sqlString,
	}).Info("OK")

	return nil
}

// Update update the given user.
func (us *UsersService) Update(u *mediumexample.User) error {
	sqlString, args, err := squirrel.Update("users").
		Set("email", u.Email).
		Set("first_name", u.FirstName).
		Set("last_name", u.LastName).
		Set("phone", u.Phone).
		Set("birthdate", u.Birthdate).
		Where("id = ?", u.ID).
		Suffix("returning *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	row := us.Store.QueryRowx(sqlString, args...)
	if err := row.StructScan(u); err != nil {
		return us.userError(sqlString, args, err)
	}

	us.Logger.WithFields(logrus.Fields{
		"query": sqlString,
	}).Info("OK")

	return nil
}

// Delete logical delete.
func (us *UsersService) Delete(u *mediumexample.User) error {
	sqlString, args, err := squirrel.Update("users").
		Set("deleted_at", time.Now()).
		Where("id = ?", u.ID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	row := us.Store.QueryRowx(sqlString, args...)
	if err := row.StructScan(u); err != nil {
		return us.userError(sqlString, args, err)
	}

	us.Logger.WithFields(logrus.Fields{
		"query": sqlString,
	}).Info("OK")

	return nil
}

func (us *UsersService) userError(sqlText string, args interface{}, err error) error {
	pqErr, ok := err.(*pq.Error)
	if !ok {
		us.Logger.WithFields(logrus.Fields{
			"query": sqlText,
		}).Error(err)
		if err == sql.ErrNoRows {
			return ErrAddressNotFound
		}

		return err
	}

	errMsg, ok := errorCodeNames[pqErr.Code]
	if !ok {
		us.Logger.WithFields(logrus.Fields{
			"query": sqlText,
		}).Error(errMsg)
		return err
	}

	us.Logger.WithFields(logrus.Fields{
		"query": sqlText,
	}).Error(pqErr)
	return errMsg
}
