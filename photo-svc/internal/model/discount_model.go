package model

import (
	"be-yourmoments/photo-svc/internal/enum"
	"time"
)

type CreateCreatorDiscountRequest struct {
	CreatorId    string            `validate:"required"`
	Name         string            `json:"name" validate:"required"`
	MinQuantity  int               `json:"min_quantity" validate:"required"`
	DiscountType enum.DiscountType `json:"discount_type" validate:"required"`
	Value        int32             `json:"value" validate:"required"`
	Active       bool              `json:"active" validate:"required"`
}

type CreatorDiscountResponse struct {
	Id           string            `json:"id"`
	CreatorId    string            `json:"creator_id"`
	Name         string            `json:"name"`
	MinQuantity  int               `json:"min_quantity"`
	DiscountType enum.DiscountType `json:"discount_type"`
	Value        int32             `json:"value"`
	Active       bool              `json:"active"`
	CreatedAt    *time.Time        `json:"created_at"`
	UpdatedAt    *time.Time        `json:"updated_at"`
}

type GetCreatorDiscountRequest struct {
	Id string `json:"id" validate:"required"`
}

type ActivateCreatorDiscountRequest struct {
	Id string `json:"id" validate:"required"`
}

type DeactivateCreatorDiscountRequest struct {
	Id string `json:"id" validate:"required"`
}
