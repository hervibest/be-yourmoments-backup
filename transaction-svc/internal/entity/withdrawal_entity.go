package entity

import "time"

type Withdrawal struct {
	Id           string     `db:"id"`
	WalletId     string     `db:"wallet_id"`
	BankWalletId string     `db:"bank_wallet_id"`
	Amount       int        `db:"amount"`
	Status       string     `db:"status"`
	Description  string     `db:"description"`
	CreatedAt    *time.Time `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}
