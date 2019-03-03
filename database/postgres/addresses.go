package postgres

import (
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	mediumexample "github.com/rodrwan/medium-example"
	"github.com/sirupsen/logrus"
)

var (
	// ErrAddressNotFound address not found.
	ErrAddressNotFound = errors.New("address not found")
)

// AddressesService is a service that query addresses table using SQLExecutor.
type AddressesService struct {
	Store  SQLExecutor
	Logger *logrus.Logger
}

// Get gets a user address by userID
func (as *AddressesService) Get(userID string) (*mediumexample.Address, error) {
	q := squirrel.Select("*").From("addresses")

	if userID != "" {
		q = q.Where("user_id = ?", userID)
	}

	sqlString, args, err := q.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var a mediumexample.Address

	row := as.Store.QueryRowx(sqlString, args...)
	if err := row.StructScan(&a); err != nil {
		return nil, as.addressesError(sqlString, args, err)
	}

	as.Logger.WithFields(logrus.Fields{
		"query": sqlString,
	}).Info("OK")

	return &a, nil
}

// Create creates the given address.
func (as *AddressesService) Create(a *mediumexample.Address) error {
	sql, args, err := squirrel.Insert("addresses").
		Columns("user_id", "address_line", "city", "locality", "region", "country", "postal_code").
		Values(a.UserID, a.AddressLine, a.City, a.Locality, a.Region, a.Country, a.PostalCode).
		Suffix("returning *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	row := as.Store.QueryRowx(sql, args...)
	if err := row.StructScan(a); err != nil {
		return as.addressesError(sql, args, err)
	}

	as.Logger.WithFields(logrus.Fields{
		"query": sql,
	}).Info("OK")

	return nil
}

// Update update the given address.
func (as *AddressesService) Update(a *mediumexample.Address) error {
	sql, args, err := squirrel.Update("addresses").
		Set("address_line", a.AddressLine).
		Set("city", a.City).
		Set("locality", a.Locality).
		Set("region", a.Region).
		Set("country", a.Country).
		Set("postal_code", a.PostalCode).
		Where("user_id = ?", a.UserID).
		Suffix("returning *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	row := as.Store.QueryRowx(sql, args...)
	if err := row.StructScan(a); err != nil {
		return as.addressesError(sql, args, err)
	}

	as.Logger.WithFields(logrus.Fields{
		"query": sql,
	}).Info("OK")

	return nil
}

// Delete logical delete.
func (as *AddressesService) Delete(a *mediumexample.Address) error {
	sql, args, err := squirrel.Delete("addresses").
		Where("user_id = ?", a.UserID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	row := as.Store.QueryRowx(sql, args...)
	if err := row.StructScan(a); err != nil {
		return as.addressesError(sql, args, err)
	}

	as.Logger.WithFields(logrus.Fields{
		"query": sql,
	}).Info("OK")

	return nil
}

func (as *AddressesService) addressesError(sqlText string, args interface{}, err error) error {
	pqErr, ok := err.(*pq.Error)
	if !ok {
		as.Logger.WithFields(logrus.Fields{
			"query": sqlText,
		}).Error(err)
		if err == sql.ErrNoRows {
			return ErrAddressNotFound
		}

		return err
	}

	errMsg, ok := errorCodeNames[pqErr.Code]
	if !ok {
		as.Logger.WithFields(logrus.Fields{
			"query": sqlText,
		}).Error(errMsg)
		return err
	}

	as.Logger.WithFields(logrus.Fields{
		"query": sqlText,
	}).Error(pqErr)
	return errMsg
}
