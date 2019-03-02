package mediumexample

import "time"

// Address hold address information returned by AddressStore
type Address struct {
	UserID      string `json:"user_id,omitempty" db:"user_id"`
	AddressLine string `json:"address_line,omitempty" db:"address_line"`
	City        string `json:"city,omitempty" db:"city"`
	Locality    string `json:"locality,omitempty" db:"locality"`
	Region      string `json:"region,omitempty" db:"region"`
	Country     string `json:"country,omitempty" db:"country"`
	PostalCode  int    `json:"postal_code,omitempty" db:"postal_code"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
