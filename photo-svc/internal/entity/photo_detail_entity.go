package entity

import (
	"time"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
)

type PhotoDetail struct {
	Id              string               `db:"id"`
	PhotoId         string               `db:"photo_id"`
	FileName        string               `db:"file_name"`
	FileKey         string               `db:"file_key"`
	Size            int64                `db:"size"`
	Type            string               `db:"type"`
	Checksum        string               `db:"checksum"`
	Width           int32                `db:"width"`
	Height          int32                `db:"height"`
	Url             string               `db:"url"`
	YourMomentsType enum.YourMomentsType `db:"your_moments_type"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
