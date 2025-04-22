package repository

type WishlistRepository interface {
	GetWishilist()
}

type wishlistRepository struct {
}

// func NewWishlistRepository() WishlistRepository {
// 	return &wishlistRepository{}
// }
