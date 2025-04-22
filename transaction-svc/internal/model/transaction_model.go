package model

import (
	"encoding/json"
	"time"
)

type CreateTransactionRequest struct {
	UserId   string   `validate:"required"`
	PhotoIds []string `json:"photo_ids" validate:"required"`
}

type CreateTransactionResponse struct {
	TransactionId string
	SnapToken     string
	RedirectURL   string
}

type UpdateTransactionWebhookRequest struct {
	TransactionType   string          `json:"transaction_type"`
	TransactionTime   string          `json:"transaction_time"`
	TransactionStatus string          `json:"transaction_status" validate:"required"`
	TransactionID     string          `json:"transaction_id"`
	StatusMessage     string          `json:"status_message"`
	StatusCode        string          `json:"status_code" validate:"required"`
	SignatureKey      string          `json:"signature_key" validate:"required"`
	ReferenceID       string          `json:"reference_id"`
	PaymentType       string          `json:"payment_type"`
	OrderID           string          `json:"order_id" validate:"required"`
	Metadata          json.RawMessage `json:"metadata"` // Bisa pakai map[string]interface{} jika ingin langsung decode
	MerchantID        string          `json:"merchant_id"`
	GrossAmount       string          `json:"gross_amount" validate:"required"`
	FraudStatus       string          `json:"fraud_status"`
	ExpiryTime        string          `json:"expiry_time"`
	Currency          string          `json:"currency"`
	Acquirer          string          `json:"acquirer"`
	Body              []byte          `json:"-"` // Untuk simpan raw body jika dibutuhkan untuk verifikasi signature, dll
}

// type CreateTransactionRequest struct {
// 	Id                       string          `json:"id" validate:"required"`
// 	Status                   string          `json:"status" validate:"required"`
// 	TransactionMethodId      string          `json:"transaction_method_id"`
// 	TransactionTypeId        string          `json:"transaction_type_id"`
// 	PaymentTypeId            string          `json:"payment_type_id"`
// 	PaymentAt                *time.Time      `json:"payment_at"`
// 	CheckoutAt               *time.Time      `json:"checkout_at"`
// 	SnapToken                string          `json:"snap_token" validate:"required"`
// 	ExternalStatus           string          `json:"external_status" `
// 	ExternalCallbackResponse json.RawMessage `json:"external_callback_response" `
// 	Amount                   int32           `json:"amount" validate:"required"`
// 	CreatedAt                time.Time       `json:"created_at" validate:"required"`
// 	UpdatedAt                time.Time       `json:"updated_at" validate:"required"`
// }

type TransactionResponse struct {
	Id                       string          `json:"id" validate:"required"`
	Status                   string          `json:"status" validate:"required"`
	TransactionMethodId      string          `json:"transaction_method_id"`
	TransactionTypeId        string          `json:"transaction_type_id"`
	PaymentTypeId            string          `json:"payment_type_id"`
	PaymentAt                *time.Time      `json:"payment_at"`
	CheckoutAt               *time.Time      `json:"checkout_at"`
	SnapToken                string          `json:"snap_token" validate:"required"`
	ExternalStatus           string          `json:"external_status" `
	ExternalCallbackResponse json.RawMessage `json:"external_callback_response" `
	Amount                   int32           `json:"amount" validate:"required"`
	CreatedAt                time.Time       `json:"created_at" validate:"required"`
	UpdatedAt                time.Time       `json:"updated_at" validate:"required"`
}

type PaymentSnapshotRequest struct {
	OrderID     string `json:"order_id,omitempty"`
	GrossAmount int64  `json:"gross_amount,omitempty"`
	Email       string `json:"email,omitempty"`
}
