package model

import "time"

type CreateReviewRequest struct {
	TransactionDetailId string  `json:"transaction_detail_id" validate:"required"`
	CreatorId           string  `json:"creator_id" validate:"required"`
	UserId              string  `json:"user_id" validate:"required"`
	Rating              int     `json:"rating" validate:"required"`
	Comment             *string `json:"comment"`
}

type CreatorReviewResponse struct {
	Id                  string
	TransactionDetailId string
	CreatorId           string
	UserId              string
	Rating              int
	Comment             *string
	CreatedAt           *time.Time
	UpdatedAt           *time.Time
}

type GetAllReviewRequest struct {
	Rating int    `json:"username" `
	Order  string `json:"order"`
	Page   int    `json:"page" validate:"required"`
	Size   int    `json:"size" validate:"required"`
}
