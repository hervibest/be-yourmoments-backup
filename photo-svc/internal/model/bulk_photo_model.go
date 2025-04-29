package model

import (
	"time"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
)

type GetBulkPhotoDetailRequest struct {
	BulkPhotoId string `validate:"required"`
	CreatorId   string `validate:"required"`
}

type GetBulkPhotoDetailResponse struct {
	Id        string               `json:"id"`
	CreatorId string               `json:"creator_id"`
	Status    enum.BulkPhotoStatus `json:"status"`
	Photo     *[]*PhotoResponse    `json:"photo"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

type PhotoResponse struct {
	Id             string    `json:"id"`
	CreatorId      string    `json:"creator_id"`
	Title          string    `json:"title"`
	OwnedByUserId  *string   `json:"owned_by_user_id"`
	CompressedUrl  *string   `json:"compressed_url"`
	IsThisYouURL   *string   `json:"is_this_you_url"`
	YourMomentsUrl *string   `json:"your_moments_url"`
	CollectionUrl  *string   `json:"collection_url"`
	Price          int32     `json:"price"`
	PriceStr       string    `json:"price_str"`
	Latitude       *float64  `json:"latitude"`
	Longitude      *float64  `json:"longitude"`
	Description    *string   `json:"description"`
	OriginalAt     time.Time `json:"original_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
