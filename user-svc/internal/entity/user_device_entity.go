package entity

import (
	"time"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/enum"
)

type UserDevice struct {
	Id        string
	UserId    string
	Token     string
	Platform  enum.PlatformTypeEnum
	CreatedAt *time.Time
}
