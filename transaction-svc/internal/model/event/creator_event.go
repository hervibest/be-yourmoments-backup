package event

import "time"

type CreatorEvent struct {
	Id        string     `json:"id" db:"id"`
	UserId    string     `json:"user_id" db:"user_id"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}
