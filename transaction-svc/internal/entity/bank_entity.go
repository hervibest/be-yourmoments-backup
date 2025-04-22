package entity

import (
	"database/sql"
	"time"
)

type Bank struct {
	Id        string         `db:"id"`
	BankCode  string         `db:"bank_code"`
	Name      string         `db:"name"`
	Alias     sql.NullString `db:"alias"`
	SwiftCode sql.NullString `db:"swift_code"`
	LogoUrl   sql.NullString `db:"logo_url"`
	CreatedAt *time.Time     `db:"created_at"`
	UpdatedAt *time.Time     `db:"updated_at"`
}
