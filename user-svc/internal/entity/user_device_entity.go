package entity

import (
	"time"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/enum"
)

type UserDevice struct {
	Id        string                `db:"id"`
	UserId    string                `db:"user_id"`
	Token     string                `db:"token"`
	Platform  enum.PlatformTypeEnum `db:"platform"`
	CreatedAt *time.Time            `db:"created_at"`
}
