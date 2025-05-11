package model

import "time"

type CreateBankWalletRequest struct {
	BankId        string `json:"bank_id" validate:"required"`
	WalletId      string `validate:"required"`
	FullName      string `json:"full_name" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
}

type DeleteBankWalletRequest struct {
	Id string `json:"id" validate:"required"`
}

type BankWalletResponse struct {
	Id            string
	WalletId      string
	BankId        string
	FullName      string
	AccountNumber string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
