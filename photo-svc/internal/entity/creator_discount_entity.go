package entity

import (
	"time"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
)

type CreatorDiscount struct {
	Id           string            `db:"id"`
	CreatorId    string            `db:"creator_id"`
	Name         string            `db:"name"`
	MinQuantity  int               `db:"min_quantity"`
	DiscountType enum.DiscountType `db:"discount_type"`
	Value        int32             `db:"value"`
	Active       bool              `db:"active"`
	CreatedAt    *time.Time        `db:"created_at"`
	UpdatedAt    *time.Time        `db:"updated_at"`
}
