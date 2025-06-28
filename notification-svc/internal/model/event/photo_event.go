package event

type BulkPhotoEvent struct {
	EventID      string           `json:"uuid"`
	UserCountMap map[string]int32 `json:"user_count_map"`
}

type SingleFacecamEvent struct {
	EventID     string `json:"uuid"`
	UserID      string `json:"user_id"`
	CountPhotos int    `json:"count_photos"`
}

type SinglePhotoEvent struct {
	EventID string   `json:"uuid"`
	UserIDs []string `json:"user_ids"`
}
