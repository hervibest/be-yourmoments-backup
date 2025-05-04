package entity

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/enum"
)

type Transaction struct {
	Id                       string                 `db:"id"`
	UserId                   string                 `db:"user_id"`
	Status                   enum.TransactionStatus `db:"status"`
	TransactionMethodId      sql.NullString         `db:"transaction_method_id"`
	TransactionTypeId        sql.NullString         `db:"transaction_type_id"`
	PaymentTypeId            sql.NullString         `db:"payment_type_id"`
	PhotoIds                 json.RawMessage        `db:"photo_ids"`
	PaymentAt                *time.Time             `db:"payment_at"`
	CheckoutAt               *time.Time             `db:"checkout_at"`
	SnapToken                sql.NullString         `db:"snap_token"`
	ExternalStatus           sql.NullString         `db:"external_status"`
	ExternalCallbackResponse *json.RawMessage       `db:"external_callback_response"`
	Amount                   int32                  `db:"amount"`
	CreatedAt                *time.Time             `db:"created_at"`
	UpdatedAt                *time.Time             `db:"updated_at"`
}

type TransactionWithDetail struct {
	TransactionId            string                 `db:"transaction_id"`
	UserId                   string                 `db:"user_id"`
	Status                   enum.TransactionStatus `db:"status"`
	TransactionMethodId      sql.NullString         `db:"transaction_method_id"`
	TransactionTypeId        sql.NullString         `db:"transaction_type_id"`
	PaymentTypeId            sql.NullString         `db:"payment_type_id"`
	PhotoIds                 json.RawMessage        `db:"photo_ids"`
	PaymentAt                *time.Time             `db:"payment_at"`
	CheckoutAt               *time.Time             `db:"checkout_at"`
	SnapToken                sql.NullString         `db:"snap_token"`
	ExternalStatus           sql.NullString         `db:"external_status"`
	ExternalCallbackResponse *json.RawMessage       `db:"external_callback_response"`
	Amount                   int32                  `db:"amount"`
	CreatedAt                *time.Time             `db:"transaction_created_at"`
	UpdatedAt                *time.Time             `db:"transaction_updated_at"`

	TranscationDetailId string `db:"transaction_detail_id"`
	CreatorId           string `db:"creator_id"`
	CreatorDiscountId   string `db:"creator_discount_id"`
	IsReviewed          bool   `db:"is_reviewed"`

	TransactionItemId string        `db:"transaction_item_id"`
	PhotoId           string        `db:"photo_id"`
	Price             int32         `db:"price"`
	Discount          sql.NullInt32 `db:"discount"`
	FinalPrice        int32         `db:"final_price"`
}
