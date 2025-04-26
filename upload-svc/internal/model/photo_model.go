package model

type CreatePhotoRequest struct {
	UserId      string   `form:"user_id" validate:"required"`
	CreatorId   string   `form:"creator_id" validate:"required"`
	PriceStr    string   `form:"price_str" validate:"required"`
	Price       int      `form:"price" validate:"required"`
	Latitude    *float64 `form:"latitude" validate:"omitempty,gte=-90,lte=90"`
	Longitude   *float64 `form:"longitude" validate:"omitempty,gte=-180,lte=180"`
	Description *string  `form:"description" validate:"omitempty,max=500"`
}
type RequestUpdateProcessedPhoto struct {
	Id                     string
	PreviewUrl             string
	PreviewWithBoundingUrl string
	UserId                 []string
}

type RequestClaimPhoto struct {
	Id     string
	UserId string
}
