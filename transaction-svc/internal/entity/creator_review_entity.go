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
	Rating              int            `db:"rating"`
	Comment             sql.NullString `db:"comment"`
	CreatedAt           *time.Time     `db:"created_at"`
	UpdatedAt           *time.Time     `db:"updated_at"`
}

type TotalReviewAndRating struct {
	CreatorId   string  `db:"creator_id"`
	TotalReview int     `db:"total_review"`
	Rating      float32 `db:"average_rating"`
}
