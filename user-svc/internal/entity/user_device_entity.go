package entity

import (
	"be-yourmoments/user-svc/internal/enum"
	"time"
)

type UserDevice struct {
	Id        string
	UserId    string
	Token     string
	Platform  enum.PlatformTypeEnum
	CreatedAt *time.Time
}
