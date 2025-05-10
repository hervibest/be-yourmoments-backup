package event

type CreatorReviewCountEvent struct {
	Id          string  `json:"id"`
	Rating      float32 `json:"rating"`
	RatingCount int     `json:"rating_count"`
}
