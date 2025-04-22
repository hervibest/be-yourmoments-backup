package entity

import "time"

type TransactionWallet struct {
	Id                  string     `db:"id"`
	WalletId            string     `db:"wallet_id"`
	TransactionDetailId string     `db:"transaction_detail_id"`
	Amount              int32      `db:"amount"`
	CreatedAt           *time.Time `db:"created_at"`
	UpdatedAt           *time.Time `db:"updated_at"`
}
