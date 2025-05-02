package route

func (r *RouteConfig) SetupExploreRoute() {
	exploreRoutes := r.App.Group("/api/explore", r.AuthMiddleware)
	// exploreRoutes.Get("/", r.ExploreController.GetUserExploreSimilar)
	exploreRoutes.Get("/", r.ExploreController.GetAllExploreSimilar)
	exploreRoutes.Get("/wishlist", r.ExploreController.GetAllUserWishlist)
	exploreRoutes.Patch("/wishlist", r.ExploreController.UserAddWishlist)
	exploreRoutes.Delete("/wishlist/delete", r.ExploreController.UserDeleteWishlist)

	exploreRoutes.Get("/favorite", r.ExploreController.GetAllUserFavorite)
	exploreRoutes.Patch("/favorite", r.ExploreController.UserAddFavorite)
	exploreRoutes.Delete("/favorite/delete", r.ExploreController.UserDeleteFavorite)

	exploreRoutes.Get("/cart", r.ExploreController.GetAllUserCart)
	exploreRoutes.Patch("/cart", r.ExploreController.UserAddCart)
	exploreRoutes.Delete("/cart/delete", r.ExploreController.UserDeleteCart)
}

func (r *RouteConfig) SetupDiscountRoute() {
	exploreRoutes := r.App.Group("/api/discount", r.AuthMiddleware)
	exploreRoutes.Get("/:discountId", r.CreatorDiscountControler.GetDiscount)
	exploreRoutes.Get("/", r.CreatorDiscountControler.GetAllDiscount)
	exploreRoutes.Post("/create", r.CreatorDiscountControler.CreateDiscount)
	exploreRoutes.Put("/activate/:discountId", r.CreatorDiscountControler.ActivateDiscount)
	exploreRoutes.Put("/deactivate/:discountId", r.CreatorDiscountControler.DeactivateDiscount)
}

func (r *RouteConfig) SetupCheckoutRoute() {
	exploreRoutes := r.App.Group("/api/checkout", r.AuthMiddleware)
	exploreRoutes.Post("/preview", r.CheckoutController.PreviewCheckout)
}
