package route

func (r *RouteConfig) SetupPhotoRoute() {
	exploreRoutes := r.App.Group("/api/photo")
	// exploreRoutes.Get("/:bulkPhotoId", r.PhotoController.GetBulkPhotoDetail)
	exploreRoutes.Get("/", r.PhotoController.GetPhotoFile)
}
