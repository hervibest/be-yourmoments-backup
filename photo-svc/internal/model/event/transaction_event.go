package event

type CancelPhotosEvent struct {
	UserId   string   `json:"user_id"`
	PhotoIds []string `json:"photo_ids"`
}

type OwnerOwnPhotosEvent struct {
	UserId   string   `json:"user_id"`
	PhotoIds []string `json:"photo_ids"`
}
