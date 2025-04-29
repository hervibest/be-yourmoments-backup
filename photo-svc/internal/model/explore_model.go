package model

import (
	"time"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
)

// type ExploreUserSimilarRequest struct {

// }

type Photo struct {
	Id             string `json:"id"`
	CreatorId      string `json:"creator_id"`
	Title          string `json:"title"`
	OwnedByUserId  string `json:"owned_by_user_id"`
	CompressedUrl  string `json:"compressed_url"`
	IsThisYouURL   string `json:"is_this_you_url"`
	YourMomentsUrl string `json:"your_moments_url"`
	CollectionUrl  string `json:"collection_url"`

	Price      int32     `json:"price"`
	PriceStr   string    `json:"price_str"`
	OriginalAt time.Time `json:"original_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserSimilarPhotoRespone struct {
	Id         string                   `json:"id"`
	PhotoId    string                   `json:"photo_id"`
	UserId     string                   `json:"user_id"`
	Similarity enum.SimilarityLevelEnum `json:"similarity"`
	IsWishlist bool                     `json:"is_wishlist"`
	IsResend   bool                     `json:"is_resend"`
	IsCart     bool                     `json:"is_cart"`
	IsFavorite bool                     `json:"is_favorite"`
	CreatedAt  time.Time                `json:"created_at"`
	UpdatedAt  time.Time                `json:"updated_at"`
}

type PhotoUrlResponse struct {
	IsThisYouURL   string `json:"is_this_you_url"`
	YourMomentsUrl string `json:"your_moments_url"`
}

type PhotoStageResponse struct {
	IsWishlist bool `json:"is_wishlist"`
	IsResend   bool `json:"is_resend"`
	IsCart     bool `json:"is_cart"`
	IsFavorite bool `json:"is_favorite"`
}

type ExploreUserSimilarResponse struct {
	PhotoId    string                   `json:"photo_id"`
	UserId     string                   `json:"user_id"`
	Similarity enum.SimilarityLevelEnum `json:"similarity"`
	PhotoStage *PhotoStageResponse      `json:"stage"`
	CreatorId  string                   `json:"creator_id"`
	Title      string                   `json:"title"`
	PhotoUrl   *PhotoUrlResponse        `json:"url"`
	Price      int32                    `json:"price"`
	PriceStr   string                   `json:"price_str"`
	Discount   *CreatorDiscountResponse `json:"discount"`
	OriginalAt time.Time                `json:"original_at"`
	CreatedAt  time.Time                `json:"created_at"`
	UpdatedAt  time.Time                `json:"updated_at"`
}

type GetAllExploreSimilarRequest struct {
	UserId string `validate:"required"`
	Page   int    `json:"page" validate:"required"`
	Size   int    `json:"size" validate:"required"`
}

type GetAllWishlistRequest struct {
	UserId string `validate:"required"`
	Page   int    `json:"page" validate:"required"`
	Size   int    `json:"size" validate:"required"`
}

type UserAddWishlistRequest struct {
	UserId  string `json:"user_id" validate:"required"`
	PhotoId string `json:"photo_id" validate:"required"`
}

type UserDeleteWishlistReqeust struct {
	UserId  string `json:"user_id" validate:"required"`
	PhotoId string `json:"photo_id" validate:"required"`
}

type GetAllFavoriteRequest struct {
	UserId string `validate:"required"`
	Page   int    `json:"page" validate:"required"`
	Size   int    `json:"size" validate:"required"`
}

type UserAddFavoriteRequest struct {
	UserId  string `json:"user_id" validate:"required"`
	PhotoId string `json:"photo_id" validate:"required"`
}

type UserDeleteFavoriteReqeust struct {
	UserId  string `json:"user_id" validate:"required"`
	PhotoId string `json:"photo_id" validate:"required"`
}

type GetAllCartRequest struct {
	UserId string `validate:"required"`
	Page   int    `json:"page" validate:"required"`
	Size   int    `json:"size" validate:"required"`
}

type UserAddCartRequest struct {
	UserId  string `json:"user_id" validate:"required"`
	PhotoId string `json:"photo_id" validate:"required"`
}

type UserDeleteCartReqeust struct {
	UserId  string `json:"user_id" validate:"required"`
	PhotoId string `json:"photo_id" validate:"required"`
}
