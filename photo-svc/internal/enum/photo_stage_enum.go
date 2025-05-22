package enum

type PhotoStageEnum string

const (
	PhotoStageWishlist PhotoStageEnum = "WISHLIST"
	PhotoStageFavorite PhotoStageEnum = "FAVORITE"
	PhotoStageCart     PhotoStageEnum = "CART"
)
