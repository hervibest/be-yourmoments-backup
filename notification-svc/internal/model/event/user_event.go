package event

import "github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper/enum"

type UserDeviceEvent struct {
	UserID      string                `json:"user_id" validate:"required"`
	DeviceToken string                `json:"device_token" validate:"required"`
	Platform    enum.PlatformTypeEnum `json:"platform" validate:"required"`
}
