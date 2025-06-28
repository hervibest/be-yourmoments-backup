package event

import (
	"time"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/enum"
)

type UserEvent struct {
	Id        string     `json:"id" db:"id"`
	Username  string     `json:"username" db:"username"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type UserDeviceEvent struct {
	UserID      string                `json:"user_id" validate:"required"`
	DeviceToken string                `json:"device_token" validate:"required"`
	Platform    enum.PlatformTypeEnum `json:"platform" validate:"required"`
}
