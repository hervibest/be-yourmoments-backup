package entity

import (
	"database/sql"
	"time"
)

type TransactionItem struct {
	Id                  string        `db:"id"`
	TransactionDetailId string        `db:"transaction_detail_id"`
	PhotoId             string        `db:"photo_id"`
	Price               int32         `db:"price"`
	Discount            sql.NullInt32 `db:"discount"`
	FinalPrice          int32         `db:"final_price"`
	CreatedAt           *time.Time    `db:"created_at"`
	UpdatedAt           *time.Time    `db:"updated_at"`
}
