package model

import (
	"time"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
)

type CreateCreatorDiscountRequest struct {
	CreatorId    string            `json:"creator_id" validate:"required"`
	Name         string            `json:"name" validate:"required"`
	MinQuantity  int               `json:"min_quantity" validate:"required"`
	DiscountType enum.DiscountType `json:"discount_type" validate:"required"`
	Value        int32             `json:"value" validate:"required"`
	IsActive     bool              `json:"is_active" validate:"required"`
}

type CreatorDiscountResponse struct {
	Id           string            `json:"id,omitempty"`
	CreatorId    string            `json:"creator_id,omitempty"`
	Name         string            `json:"name"`
	MinQuantity  int               `json:"min_quantity"`
	DiscountType enum.DiscountType `json:"discount_type"`
	Value        int32             `json:"value"`
	IsActive     bool              `json:"is_active"`
	CreatedAt    *time.Time        `json:"created_at,omitempty"`
	UpdatedAt    *time.Time        `json:"updated_at,omitempty"`
}

type GetCreatorDiscountRequest struct {
	Id        string `json:"id" validate:"required"`
	CreatorId string `json:"creator_id" validate:"required"`
}

type ActivateCreatorDiscountRequest struct {
	Id        string `json:"id" validate:"required"`
	CreatorId string `json:"creator_id" validate:"required"`
}

type DeactivateCreatorDiscountRequest struct {
	Id        string `json:"id" validate:"required"`
	CreatorId string `json:"creator_id" validate:"required"`
}
