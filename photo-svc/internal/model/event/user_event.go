package event

import "time"

type UserEvent struct {
	Id        string     `json:"id" db:"id"`
	Username  string     `json:"username" db:"username"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}
