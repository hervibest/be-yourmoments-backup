package entity

import (
	"time"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/enum"
)

type BulkPhoto struct {
	Id              string
	CreatorId       string
	BulkPhotoStatus enum.BulkPhotoStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
