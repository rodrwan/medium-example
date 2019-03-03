package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	mediumexample "github.com/rodrwan/medium-example"
)

var (
	// ErrAddressNotFound address not found.
	ErrAddressNotFound = errors.New("address not found")
)

// AddressesService is a service that query addresses table using SQLExecutor.
type AddressesService struct {
	Store  SQLExecutorContext
	Logger Logger
}

// Get gets a user address by userID
func (as *AddressesService) Get(ctx context.Context, userID string) (*mediumexample.Address, error) {
	q := squirrel.Select("*").From("addresses")

	if userID != "" {
		q = q.Where("user_id = ?", userID)
	}

	sqlString, args, err := q.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var a mediumexample.Address
	as.Logger.Print(sqlString, args...)
	row := as.Store.QueryRowxContext(ctx, sqlString, args...)
	if err := row.StructScan(&a); err != nil {
		return nil, as.addressesError(sqlString, args, err)
	}

	return &a, nil
}

// Create creates the given address.
func (as *AddressesService) Create(ctx context.Context, a *mediumexample.Address) error {
	sql, args, err := squirrel.Insert("addresses").
		Columns("user_id", "address_line", "city", "locality", "region", "country", "postal_code").
		Values(a.UserID, a.AddressLine, a.City, a.Locality, a.Region, a.Country, a.PostalCode).
		Suffix("returning *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	as.Logger.Print(sql, args...)
	row := as.Store.QueryRowxContext(ctx, sql, args...)
	if err := row.StructScan(a); err != nil {
		return as.addressesError(sql, args, err)
	}

	return nil
}

// Update update the given address.
func (as *AddressesService) Update(ctx context.Context, a *mediumexample.Address) error {
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

	as.Logger.Print(sql, args...)
	row := as.Store.QueryRowxContext(ctx, sql, args...)
	if err := row.StructScan(a); err != nil {
		return as.addressesError(sql, args, err)
	}

	return nil
}

// Delete logical delete.
func (as *AddressesService) Delete(ctx context.Context, a *mediumexample.Address) error {
	sql, args, err := squirrel.Delete("addresses").
		Where("user_id = ?", a.UserID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	as.Logger.Print(sql, args...)
	row := as.Store.QueryRowxContext(ctx, sql, args...)
	if err := row.StructScan(a); err != nil {
		return as.addressesError(sql, args, err)
	}

	return nil
}

func (as *AddressesService) addressesError(sqlText string, args interface{}, err error) error {
	pqErr, ok := err.(*pq.Error)
	if !ok {
		as.Logger.Warn(sqlText, args, err)
		if err == sql.ErrNoRows {
			return ErrAddressNotFound
		}

		return err
	}

	errMsg, ok := errorCodeNames[pqErr.Code]
	if !ok {
		as.Logger.Warn(sqlText, args, errMsg)
		return err
	}

	as.Logger.Warn(sqlText, args, pqErr)
	return errMsg
}
