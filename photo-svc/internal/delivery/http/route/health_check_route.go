package route

func (r *RouteConfig) SetupHealtCheckRoute() {
	exploreRoutes := r.App.Group("/api")
	// exploreRoutes.Get("/", r.ExploreController.GetUserExploreSimilar)
	exploreRoutes.Get("/", r.HealthCheckController.HealthCheck)
}
