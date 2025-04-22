package model

import (
	"time"
)

type CreateBankRequest struct {
	BankCode  string  `json:"bank_code" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	Alias     *string `json:"alias" validate:""`
	SwiftCode *string `json:"swift_code" validate:""`
	LogoUrl   *string `json:"logo_url" validate:""`
}

type UpdateBankRequest struct {
	BankCode  string  `json:"bank_code" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	Alias     *string `json:"alias" validate:""`
	SwiftCode *string `json:"swift_code" validate:""`
	LogoUrl   *string `json:"logo_url" validate:""`
}

type DeleteBankRequest struct {
	Id string `query:"id" json:"id"`
}

type FindBankByIdRequest struct {
	Id string `json:"id" validate:"required"`
}

type BankResponse struct {
	Id        string     `json:"id"`
	BankCode  string     `json:"bank_code"`
	Name      string     `json:"name"`
	Alias     *string    `json:"alias,omitempty"`
	SwiftCode *string    `json:"swift_code,omitempty"`
	LogoUrl   *string    `json:"logo_url,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
