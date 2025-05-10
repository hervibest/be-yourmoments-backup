package model

import "time"

type CreateCreatorRequest struct {
	UserId string `json:"user_id" validate:"required,max=100"`
}

type GetCreatorRequest struct {
	UserId string `json:"user_id" validate:"required,max=100"`
}

type UpdateCreatorTotalRatingRequest struct {
	Id          string  `json:"id"`
	Rating      float32 `json:"rating"`
	RatingCount int     `json:"rating_count"`
}

type CreatorResponse struct {
	Id          string     `json:"id"`
	UserId      string     `json:"user_id"`
	Rating      float32    `json:"rating,omitempty"`
	RatingCount int        `json:"rating_count,omitempty"`
	VerifiedAt  *time.Time `json:"verified_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
