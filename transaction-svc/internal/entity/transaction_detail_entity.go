package entity

import (
	"time"
)

type TransactionDetail struct {
	Id                string     `db:"id"`
	TransactionId     string     `db:"transaction_id"`
	CreatorId         string     `db:"creator_id"`
	SubTotalPrice     int32      `db:"subtotal_price"`
	CreatorDiscountId string     `db:"creator_discount_id"`
	IsReviewed        bool       `db:"is_reviewed"`
	CreatedAt         *time.Time `db:"created_at"`
	UpdatedAt         *time.Time `db:"updated_at"`
}
