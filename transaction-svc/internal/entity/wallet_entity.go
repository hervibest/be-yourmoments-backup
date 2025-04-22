package entity

import "time"

type Wallet struct {
	Id        string     `db:"id"`
	CreatorId string     `db:"creator_id"`
	Balance   int32      `db:"balance"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
