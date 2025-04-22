package entity

import (
	"time"
)

type BankWallet struct {
	Id            string     `db:"id"`
	WalletId      string     `db:"wallet_id"`
	BankId        string     `db:"bank_id"`
	FullName      string     `db:"full_name"`
	AccountNumber string     `db:"account_number"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}
