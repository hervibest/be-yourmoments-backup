package model

import (
	"time"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
)

type CheckoutItem struct {
	PhotoId             string            `json:"photo_id"`
	CreatorId           string            `json:"creator_id"`
	Title               string            `json:"title"`
	YourMomentsUrl      string            `json:"your_moments_url"`
	Price               int32             `json:"price"`
	Discount            int32             `json:"discount"`
	DiscountMinQuantity int               `json:"discount_min_quantity"`
	DiscountValue       int32             `json:"discount_value"`
	DiscountId          string            `json:"discount_id"`
	DiscountType        enum.DiscountType `json:"discount_type"`
	FinalPrice          int32             `json:"final_price"`
}

type PreviewCheckoutRequest struct {
	UserId   string   `json:"user_id" validate:"required"`
	PhotoIds []string `json:"photo_ids" validate:"required"`
}

type PreviewCheckoutResponse struct {
	Items         *[]*CheckoutItem `json:"items"`
	TotalPrice    int32            `json:"total_price"`
	TotalDiscount int32            `json:"total_discount"`
	CreatedAt     *time.Time       `json:"created_at"`
}

type Total struct {
	Price    int32
	Discount int32
}
