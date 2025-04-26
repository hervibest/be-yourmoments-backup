package entity

import (
	"database/sql"
	"time"
)

type Photo struct {
	Id             string          `db:"id"`
	CreatorId      string          `db:"creator_id"`
	BulkPhotoId    sql.NullString  `db:"bulk_photo_id"`
	Title          string          `db:"title"`
	OwnedByUserId  sql.NullString  `db:"owned_by_user_id"`
	CompressedUrl  sql.NullString  `db:"compressed_url"`
	IsThisYouURL   sql.NullString  `db:"is_this_you_url"`
	YourMomentsUrl sql.NullString  `db:"your_moments_url"`
	CollectionUrl  sql.NullString  `db:"collection_url"`
	Price          int32           `db:"price"`
	PriceStr       string          `db:"price_str"`
	Latitude       sql.NullFloat64 `db:"latitude"`
	Longitude      sql.NullFloat64 `db:"longitude"`
	Description    sql.NullString  `db:"description"`
	OriginalAt     time.Time       `db:"original_at"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
}
