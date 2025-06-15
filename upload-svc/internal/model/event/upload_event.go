package event

// Message structures
type ProcessPhotoMessage struct {
	PhotoID          string `json:"photo_id"`
	CreatorID        string `json:"creator_id"`
	URL              string `json:"url"`
	OriginalFilename string `json:"original_filename"`
	Timestamp        int64  `json:"timestamp"`
}

type ProcessFacecamMessage struct {
	UserID    string `json:"user_id"`
	CreatorID string `json:"creator_id"`
	URL       string `json:"url"`
	Timestamp int64  `json:"timestamp"`
}

type ProcessBulkPhotoMessage struct {
	BulkPhotoID string           `json:"bulk_photo_id"`
	CreatorID   string           `json:"creator_id"`
	Photos      []*BulkPhotoItem `json:"photos"`
	Timestamp   int64            `json:"timestamp"`
}

type BulkPhotoItem struct {
	ID               string `json:"id"`
	CollectionURL    string `json:"collection_url"`
	OriginalFilename string `json:"original_filename"`
}
