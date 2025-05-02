package entity

import (
	"database/sql"
	"time"
)

type User struct {
	Id                    string         `json:"id" db:"id"`
	Username              string         `json:"username" db:"username"`
	Email                 sql.NullString `json:"email" db:"email"`
	EmailVerifiedAt       *time.Time     `json:"email_verified_at" db:"email_verified_at"`
	Password              sql.NullString `json:"password" db:"password"`
	PhoneNumber           sql.NullString `json:"phone_number" db:"phone_number"`
	PhoneNumberVerifiedAt *time.Time     `json:"phone_number_verified_at" db:"phone_number_verified_at"`
	GoogleId              sql.NullString `json:"google_id" db:"google_id"`
	CreatedAt             *time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt             *time.Time     `json:"updated_at" db:"updated_at"`
}

func (u *User) HasEmail() bool {
	return u.Email.Valid
}

func (u *User) HasVerifiedEmail() bool {
	return u.EmailVerifiedAt != nil
}

func (u *User) HasPhoneNumber() bool {
	return u.PhoneNumber.Valid
}

func (u *User) HasVerifiedPhoneNumber() bool {
	return u.PhoneNumberVerifiedAt != nil
}

type UserPublicChat struct {
	UserId   string         `db:"user_id"`
	Username string         `db:"username"`
	FileKey  sql.NullString `db:"file_key"`
}

type UserDetail struct {
	Id          string         `json:"id" db:"id"`
	Username    string         `json:"username" db:"username"`
	Email       sql.NullString `json:"email" db:"email"`
	PhoneNumber sql.NullString `json:"phone_number" db:"phone_number"`
	Similarity  uint32         `json:"similarity" db:"similarity"`
}
