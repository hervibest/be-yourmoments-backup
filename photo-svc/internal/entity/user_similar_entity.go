package entity

import (
	"database/sql"
	"time"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
)

type UserSimilarPhoto struct {
	PhotoId    string                   `db:"photo_id"`
	UserId     string                   `db:"user_id"`
	Similarity enum.SimilarityLevelEnum `db:"similarity"`
	IsWishlist bool                     `db:"is_wishlist"`
	IsResend   bool                     `db:"is_resend"`
	IsCart     bool                     `db:"is_cart"`
	IsFavorite bool                     `db:"is_favorite"`
	CreatedAt  time.Time                `db:"created_at"`
	UpdatedAt  time.Time                `db:"updated_at"`
}

type Explore struct {
	PhotoId    string                   `db:"photo_id"`
	UserId     string                   `db:"user_id"`
	Similarity enum.SimilarityLevelEnum `db:"similarity"`
	IsWishlist bool                     `db:"is_wishlist"`
	IsResend   bool                     `db:"is_resend"`
	IsCart     bool                     `db:"is_cart"`
	IsFavorite bool                     `db:"is_favorite"`
	CreatorId  string                   `db:"creator_id"`
	Title      string                   `db:"title"`

	IsThisYouURL   sql.NullString `db:"is_this_you_url"`
	YourMomentsUrl sql.NullString `db:"your_moments_url"`
	Price          int32          `db:"price"`
	PriceStr       string         `db:"price_str"`
	OriginalAt     time.Time      `db:"original_at"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`

	Name         sql.NullString `db:"name"`
	MinQuantity  sql.NullInt32  `db:"min_quantity"`
	DiscountType sql.NullString `db:"discount_type"`
	Value        sql.NullInt32  `db:"value"`
	Active       sql.NullBool   `db:"active"`
}
