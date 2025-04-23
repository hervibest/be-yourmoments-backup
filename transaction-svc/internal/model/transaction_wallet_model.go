package model

import "time"

type GetAllTransactionWallet struct {
	WalletId string `json:"wallet_id" validate:"required"`
	Max      string `json:"max"`
	Min      string `json:"min"`
	Order    string `json:"time_order" validate:"required"`
	Page     int    `json:"page" validate:"required"`
	Size     int    `json:"size" validate:"required"`
}

type TransactionWalletResponse struct {
	Id                  string     `json:"id"`
	WalletId            string     `json:"wallet_id"`
	TransactionDetailId string     `json:"transaction_detail_id"`
	Amount              int32      `json:"amount"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
}
