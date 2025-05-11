package model

import "time"

type CreateWithdrawalRequest struct {
	WalletId     string `json:"wallet_id" validate:"required"`
	BankWalletId string `json:"bank_wallet_id" validate:"required"`
	Amount       int    `json:"amount" validate:"required"`
}

type UpdateWithdrawalStatusRequest struct {
	Id     string `json:"id" validate:"required"`
	Status string `json:"status" validate:"required"`
}

type FindWithdrawalById struct {
	Id string `json:"id" validate:"required"`
}

type WithdrawalResponse struct {
	Id           string     `json:"id"`
	WalletId     string     `json:"wallet_id"`
	BankWalletId string     `json:"bank_wallet_id"`
	Amount       int        `json:"amount"`
	Status       string     `json:"status"`
	Description  string     `json:"description"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
