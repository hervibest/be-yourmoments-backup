package route

func (r *RouteConfig) SetupPhotoRoute() {
	exploreRoutes := r.App.Group("/api/photo", r.AuthMiddleware)
	exploreRoutes.Get("/:bulkPhotoId", r.PhotoController.GetBulkPhotoDetail)
}
