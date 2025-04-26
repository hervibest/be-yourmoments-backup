package entity

import (
	"be-yourmoments/photo-svc/internal/enum"
	"database/sql"
	"time"
)

type BulkPhoto struct {
	Id              string               `db:"id"`
	CreatorId       string               `db:"creator_id"`
	BulkPhotoStatus enum.BulkPhotoStatus `db:"bulk_photo_status"`
	CreatedAt       time.Time            `db:"created_at"`
	UpdatedAt       time.Time            `db:"updated_at"`
}

type BulkPhotoDetail struct {
	BulkPhotoId        string               `db:"bulk_photo_id"`
	BulkPhotoCreatorId string               `db:"bulk_photo_creator_id"`
	BulkPhotoStatus    enum.BulkPhotoStatus `db:"bulk_photo_status"`
	BulkPhotoCreatedAt time.Time            `db:"bulk_photo_created_at"`
	BulkPhotoUpdatedAt time.Time            `db:"bulk_photo_updated_at"`

	PhotoId             string          `db:"photo_id"`
	PhotoCreatorId      string          `db:"photo_creator_id"`
	PhotoTitle          string          `db:"photo_title"`
	PhotoOwnedByUserId  sql.NullString  `db:"photo_owned_by_user_id"`
	PhotoCompressedUrl  sql.NullString  `db:"photo_compressed_url"`
	PhotoIsThisYouURL   sql.NullString  `db:"photo_is_this_you_url"`
	PhotoYourMomentsUrl sql.NullString  `db:"photo_your_moments_url"`
	PhotoCollectionUrl  sql.NullString  `db:"photo_collection_url"`
	PhotoPrice          int32           `db:"photo_price"`
	PhotoPriceStr       string          `db:"photo_price_str"`
	PhotoLatitude       sql.NullFloat64 `db:"photo_latitude"`
	PhotoLongitude      sql.NullFloat64 `db:"photo_longitude"`
	PhotoDescription    sql.NullString  `db:"photo_description"`
	PhotoOriginalAt     time.Time       `db:"photo_original_at"`
	PhotoCreatedAt      time.Time       `db:"photo_created_at"`
	PhotoUpdatedAt      time.Time       `db:"photo_updated_at"`
}
