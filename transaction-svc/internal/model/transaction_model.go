package model

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/enum"
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
	SettlementTime    string          `json:"settlement_time" validate:"required"`
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

type GetTransactionWithDetail struct {
	TransactionId string `validate:"required"`
	UserID        string `validate:"required"`
}

// type TransactionWithDetail struct {
// 	TransactionId       string                 `json:"transaction_id"`
// 	UserId              string                 `json:"user_id"`
// 	Status              enum.TransactionStatus `json:"status"`
// 	TransactionMethodId *string                `json:"transaction_method_id,omitempty"`
// 	TransactionTypeId   *string                `json:"transaction_type_id,omitempty"`
// 	PaymentTypeId       *string                `json:"payment_type_id,omitempty"`
// 	PaymentAt           *time.Time             `json:"payment_at,omitempty"`
// 	CheckoutAt          *time.Time             `json:"checkout_at,omitempty"`
// 	Amount              int32                  `json:"amount"`
// 	CreatedAt           *time.Time             `json:"transaction_created_at"`
// 	UpdatedAt           *time.Time             `json:"transaction_updated_at"`

// 	CreatorId         string `json:"creator_id"`
// 	CreatorDiscountId string `json:"creator_discount_id"`
// 	IsReviewed        bool   `json:"is_reviewed"`

// 	PhotoId    string        `json:"photo_id"`
// 	Price      int32         `json:"price"`
// 	Discount   sql.NullInt32 `json:"discount"`
// 	FinalPrice int32         `json:"final_price"`
// 	Url        string        `json:"url"`

// 	Title           string    `json:"title"`
// 	Latitude        *float64  `json:"latitude,omitempty"`
// 	Longitude       *float64  `json:"longitude,omitempty"`
// 	Description     *string   `json:"description,omitempty"`
// 	PhotoOriginalAt time.Time `json:"photo_original_at"`
// 	PhotoCreatedAt  time.Time `json:"photo_created_at"`
// 	PhotoUpdatedAt  time.Time `json:"photo_updated_at"`

// 	FileName        string               `json:"file_name"`
// 	Size            int64                `json:"size"`
// 	Type            string               `json:"type"`
// 	Width           int32                `json:"width"`
// 	Height          int32                `json:"height"`
// 	YourMomentsType enum.YourMomentsType `json:"your_moments_type"`
// }

type TransactionDetailResponse struct {
	CreatorId         string            `json:"creator_id"`
	CreatorDiscountId string            `json:"creator_discount_id"`
	IsReviewed        bool              `json:"is_reviewed"`
	Photo             *[]*PhotoResponse `json:"photos"`
}

type PhotoResponse struct {
	PhotoId    string        `json:"photo_id"`
	Price      int32         `json:"price"`
	Discount   sql.NullInt32 `json:"discount"`
	FinalPrice int32         `json:"final_price"`
	Url        string        `json:"url"`

	Title           string    `json:"title"`
	Latitude        *float64  `json:"latitude,omitempty"`
	Longitude       *float64  `json:"longitude,omitempty"`
	Description     *string   `json:"description,omitempty"`
	PhotoOriginalAt time.Time `json:"photo_original_at"`
	PhotoCreatedAt  time.Time `json:"photo_created_at"`
	PhotoUpdatedAt  time.Time `json:"photo_updated_at"`

	FileName        string               `json:"file_name"`
	Size            int64                `json:"size"`
	Type            string               `json:"type"`
	Width           int32                `json:"width"`
	Height          int32                `json:"height"`
	YourMomentsType enum.YourMomentsType `json:"your_moments_type"`
}

type TransactionWithDetail struct {
	TransactionId       string                        `json:"transaction_id"`
	UserId              string                        `json:"user_id"`
	Status              enum.TransactionStatus        `json:"status"`
	TransactionMethodId *string                       `json:"transaction_method_id,omitempty"`
	TransactionTypeId   *string                       `json:"transaction_type_id,omitempty"`
	PaymentTypeId       *string                       `json:"payment_type_id,omitempty"`
	PaymentAt           *time.Time                    `json:"payment_at,omitempty"`
	CheckoutAt          *time.Time                    `json:"checkout_at,omitempty"`
	Amount              int32                         `json:"amount"`
	CreatedAt           *time.Time                    `json:"transaction_created_at"`
	UpdatedAt           *time.Time                    `json:"transaction_updated_at"`
	TransactionDetail   *[]*TransactionDetailResponse `json:"transaction_detail_response"`
}

type UserTransaction struct {
	Id                  string                 `json:"id"`
	UserId              string                 `json:"user_id"`
	Status              enum.TransactionStatus `json:"status"`
	TransactionMethodId *string                `json:"transaction_method_id,omitempty"`
	TransactionTypeId   *string                `json:"transaction_type_id,omitempty"`
	PaymentTypeId       *string                `json:"payment_type_id,omitempty"`
	PaymentAt           *time.Time             `json:"payment_at,omitempty"`
	CheckoutAt          *time.Time             `json:"checkout_at,omitempty"`
	Amount              int32                  `json:"amount"`
	// CreatedAt           *time.Time             `json:"created_at"`
	// UpdatedAt           *time.Time             `json:"updated_at"`
}

type GetAllUsertTransaction struct {
	UserId string `json:"user_id" validate:"required"`
	Order  string `json:"order" validate:"required"`
	Page   int    `json:"page" validate:"required"`
	Size   int    `json:"size" validate:"required"`
}
