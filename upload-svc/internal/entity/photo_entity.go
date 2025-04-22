package entity

import (
	"database/sql"
	"time"
)

type Photo struct {
	Id             string `db:"id"`
	UserId         string
	CreatorId      string          `db:"creator_id"`
	Title          string          `db:"title"`
	OwnedByUserId  string          `db:"owned_by_user_id"`
	CompressedUrl  string          `db:"compressed_url"`
	IsThisYouURL   string          `db:"is_this_you_url"`
	YourMomentsUrl string          `db:"your_moments_url"`
	CollectionUrl  string          `db:"collection_url"`
	Latitude       sql.NullFloat64 `db:"latitude"`
	Longitude      sql.NullFloat64 `db:"longitude"`
	Description    sql.NullString  `db:"description"`

	Price      int       `db:"price"`
	PriceStr   string    `db:"price_str"`
	OriginalAt time.Time `db:"original_at"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
