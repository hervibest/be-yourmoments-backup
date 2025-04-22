package entity

import "time"

type Creator struct {
	Id          string     `db:"id"`
	UserId      string     `db:"user_id"`
	Rating      float32    `db:"rating"`
	RatingCount int        `db:"rating_count"`
	VerifiedAt  *time.Time `db:"verified_at"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}
