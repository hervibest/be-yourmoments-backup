package entity

import (
	"be-yourmoments/upload-svc/internal/enum"
	"time"
)

type BulkPhoto struct {
	Id              string
	CreatorId       string
	BulkPhotoStatus enum.BulkPhotoStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
