package entity

import (
	"database/sql"
	"time"
)

type CreatorReview struct {
	Id                  string         `db:"id"`
	TransactionDetailId string         `db:"transaction_detail_id"`
	CreatorId           string         `db:"creator_id"`
	UserId              string         `db:"user_id"`
	Star                int            `db:"star"`
	Comment             sql.NullString `db:"comment"`
	CreatedAt           *time.Time     `db:"created_at"`
	UpdatedAt           *time.Time     `db:"updated_at"`
}
