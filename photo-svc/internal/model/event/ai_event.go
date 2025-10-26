package event

import "time"

type UserSimilarEvent struct {
	PhotoDetail      PhotoDetail        `json:"photo_detail"`
	UserSimilarPhoto []UserSimilarPhoto `json:"user_similar_photo"`
}

type UserSimilarPhoto struct {
	Id         string     `json:"id"`
	PhotoID    string     `json:"photo_id"`
	UserID     string     `json:"user_id"`
	Similarity uint32     `json:"similarity"`
	IsWishlist bool       `json:"is_wishlist"`
	IsResend   bool       `json:"is_resend"`
	IsCart     bool       `json:"is_cart"`
	IsFavorite bool       `json:"is_favorite"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type BulkUserSimilarPhoto struct {
	PhotoDetail      PhotoDetail        `json:"photo_detail"`
	UserSimilarPhoto []UserSimilarPhoto `json:"user_similar_photo"`
}

type BulkUserSimilarPhotoEvent struct {
	BulkPhoto            BulkPhoto              `json:"bulk_photo"`
	BulkUserSimilarPhoto []BulkUserSimilarPhoto `json:"bulk_user_similar_photo"`
}

type UserSimiliarFacecamEvent struct {
	Facecam          Facecam            `json:"facecam"`
	UserSimilarPhoto []UserSimilarPhoto `json:"user_similar_photo"`
}
