package route

func (r *route) setupReviewRoute() {
	reviewRoute := r.app.Group("/api/review", r.authMiddleware)
	reviewRoute.Post("/create", r.reviewController.CreateReview)
	reviewRoute.Get("/", r.reviewController.GetAllReview)
}
