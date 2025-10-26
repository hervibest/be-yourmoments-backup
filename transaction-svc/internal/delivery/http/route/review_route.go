package route

func (r *route) setupReviewRoute() {
	reviewRoute := r.app.Group("/api/review", r.authMiddleware, r.creatorMiddleware, r.walletMiddleware)
	reviewRoute.Post("/create", r.reviewController.UserCreateReview)
	reviewRoute.Get("/", r.creatorMiddleware, r.reviewController.CreatorGetReview)
}
