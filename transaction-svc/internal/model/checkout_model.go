package model

type CheckoutItem struct {
	PhotoId             string `json:"photo_id"`
	CreatorId           string `json:"creator_id"`
	Title               string `json:"title"`
	YourMomentsUrl      string `json:"your_moments_url"`
	Price               int32  `json:"price"`
	Discount            int32  `json:"discount"`
	DiscountMinQuantity int    `json:"discount_min_quantity"`
	DiscountValue       int32  `json:"discount_value"`
	DiscountId          string `json:"discount_id"`
	DiscountType        string `json:"discount_type"`
	FinalPrice          int32  `json:"final_price"`
}

type Total struct {
	Price    int32
	Discount int32
}

type CheckoutItemWeb struct {
	PhotoId    string        `json:"photo_id" validate:"required"`
	CreatorId  string        `json:"creator_id" validate:"required"`
	Title      string        `json:"title" validate:"required"`
	Price      int32         `json:"price" validate:"required,gt=0"`
	Discount   *DiscountItem `json:"discount,omitempty"`
	FinalPrice int32         `json:"final_price" validate:"required,gt=0"`
}

type DiscountItem struct {
	Discount            int32  `json:"discount" validate:"required"`
	DiscountMinQuantity int    `json:"discount_min_quantity" validate:"required,gte=0"`
	DiscountValue       int32  `json:"discount_value" validate:"required,gt=0"`
	DiscountId          string `json:"discount_id" validate:"required"`
	DiscountType        string `json:"discount_type" validate:"required"`
}

type PreviewCheckoutResponse struct {
	Items         []CheckoutItemWeb `json:"items" validate:"required"`
	TotalPrice    int32             `json:"total_price" validate:"required,gt=0"`
	TotalDiscount int32             `json:"total_discount" validate:"required,gt=0"`
}

type CreateTransactionV2Request struct {
	UserId        string            `validate:"required"`
	CreatorId     string            `validate:"required"`
	Items         []CheckoutItemWeb `json:"items" validate:"required,dive"`
	TotalPrice    int32             `json:"total_price" validate:"required,gte=0"`
	TotalDiscount int32             `json:"total_discount" `
}
