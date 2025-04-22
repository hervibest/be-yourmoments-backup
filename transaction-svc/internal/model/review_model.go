package model

import "time"

type CreateReviewRequest struct {
	TransactionDetailId string
	CreatorId           string
	UserId              string
	Star                int
	Comment             *string
}

type CreatorReviewResponse struct {
	Id                  string
	TransactionDetailId string
	CreatorId           string
	UserId              string
	Star                int
	Comment             *string
	CreatedAt           *time.Time
	UpdatedAt           *time.Time
}

type GetAllReviewRequest struct {
	Star  int    `json:"username" `
	Order string `json:"order"`
	Page  int    `json:"page" validate:"required"`
	Size  int    `json:"size" validate:"required"`
}
