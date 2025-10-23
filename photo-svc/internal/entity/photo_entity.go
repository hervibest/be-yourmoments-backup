package entity

import (
	"database/sql"
	"time"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
)

type Photo struct {
	Id               string          `db:"id"`
	CreatorId        string          `db:"creator_id"`
	BulkPhotoId      sql.NullString  `db:"bulk_photo_id"`
	Title            string          `db:"title"`
	OwnedByUserId    sql.NullString  `db:"owned_by_user_id"`
	CompressedUrl    sql.NullString  `db:"compressed_url"`
	IsThisYouURL     sql.NullString  `db:"is_this_you_url"`
	YourMomentsUrl   sql.NullString  `db:"your_moments_url"`
	CollectionUrl    sql.NullString  `db:"collection_url"`
	Price            int32           `db:"price"`
	PriceStr         string          `db:"price_str"`
	Latitude         sql.NullFloat64 `db:"latitude"`
	Longitude        sql.NullFloat64 `db:"longitude"`
	Description      sql.NullString  `db:"description"`
	TotalUserSimilar int             `db:"total_user_similar"`
	OriginalAt       time.Time       `db:"original_at"`
	CreatedAt        time.Time       `db:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at"`

	FileName  sql.NullString `db:"file_name"`
	FileKey   sql.NullString `db:"file_key"`
	PhotoType sql.NullString `db:"your_moments_type"`

	Status enum.PhotoStatusEnum `db:"status"`
}

type PhotoWithDetail struct {
	Id          string          `db:"photo_id"`
	CreatorId   string          `db:"creator_id"`
	Title       string          `db:"title"`
	Price       int32           `db:"price"`
	PriceStr    string          `db:"price_str"`
	Latitude    sql.NullFloat64 `db:"latitude"`
	Longitude   sql.NullFloat64 `db:"longitude"`
	Description sql.NullString  `db:"description"`
	OriginalAt  time.Time       `db:"original_at"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`

	FileName        string               `db:"file_name"`
	FileKey         string               `db:"file_key"`
	Size            int64                `db:"size"`
	Type            string               `db:"type"`
	Width           int32                `db:"width"`
	Height          int32                `db:"height"`
	YourMomentsType enum.YourMomentsType `db:"your_moments_type"`

	Status enum.PhotoStatusEnum `db:"status"`
}
