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
