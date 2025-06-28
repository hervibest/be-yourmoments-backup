package event

import "time"

type UserSimilarEvent struct {
	PhotoDetail      PhotoDetail        `json:"photo_detail"`
	UserSimilarPhoto []UserSimilarPhoto `json:"user_similar_photo"`
}

type PhotoDetail struct {
	Id              string     `json:"id"`
	PhotoID         string     `json:"photo_id"`
	FileName        string     `json:"file_name"`
	FileKey         string     `json:"file_key"`
	Size            int64      `json:"size"`
	Type            string     `json:"type"`
	Checksum        string     `json:"checksum"`
	Width           int64      `json:"width"`
	Height          int64      `json:"height"`
	Url             string     `json:"url"`
	YourMomentsType string     `json:"your_moments_type"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type UserSimilarPhoto struct {
	Id         string     `json:"id"`
	PhotoID    string     `json:"photo_d"`
	UserID     string     `json:"file_name"`
	Similarity uint32     `json:"similarity"`
	IsWishlist bool       `json:"is_wishlist"`
	IsResend   bool       `json:"is_resend"`
	IsCart     bool       `json:"is_cart"`
	IsFavorite bool       `json:"is_favorite"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type BulkPhoto struct {
	Id              string     `json:"id"`
	CreatorId       string     `json:"creator_id"`
	BulkPhotoStatus string     `json:"bulk_photo_status"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type BulkUserSimilarPhoto struct {
	PhotoDetail      PhotoDetail        `json:"photo_detail"`
	UserSimilarPhoto []UserSimilarPhoto `json:"user_similar_photo"`
}

type BulkUserSimilarPhotoEvent struct {
	BulkPhoto            BulkPhoto              `json:"bulk_photo"`
	BulkUserSimilarPhoto []BulkUserSimilarPhoto `json:"bulk_user_similar_photo"`
}

type Facecam struct {
	Id          string     `json:"id"`
	UserId      string     `json:"user_id"`
	FileName    string     `json:"file_name"`
	FileKey     string     `json:"file_key"`
	Title       string     `json:"title"`
	Size        int64      `json:"size"`
	Type        string     `json:"type"`
	Checksum    string     `json:"checksum"`
	Url         string     `json:"url"`
	IsProcessed bool       `json:"is_processed"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type UserSimiliarFacecamEvent struct {
	Facecam          Facecam            `json:"facecam"`
	UserSimilarPhoto []UserSimilarPhoto `json:"user_similar_photos"`
}
