package event

import "time"

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
	OriginalAt  *time.Time `json:"original_at"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type BulkPhoto struct {
	Id              string     `json:"id"`
	CreatorId       string     `json:"creator_id"`
	BulkPhotoStatus string     `json:"bulk_photo_status"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type Photo struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	CreatorID      string `json:"creator"`
	Title          string `json:"title"`
	OwnedByUserID  string `json:"owned_by_user_id"`
	ComporessedURL string `json:"compressed_url"`
	IsThisYouURL   string `json:"is_this_you_url"`
	YourMomentsURL string `json:"your_moments_url"`
	CollectionURL  string `json:"collection_url"`

	Price      int32      `json:"price"`
	PriceStr   string     `json:"price_str"`
	OriginalAt *time.Time `json:"original_at"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`

	Url         string   `json:"url"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	Description *string  `json:"description"`
	BulkPhotoID *string  `json:"bulk_photo_id"`

	PhotoDetail PhotoDetail `json:"photo_detail"`
}

type PhotoDetail struct {
	Id              string     `json:"id"`
	PhotoId         string     `json:"photo_id"`
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
