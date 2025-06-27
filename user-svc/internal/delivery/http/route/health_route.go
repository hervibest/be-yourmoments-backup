package route

func (c *RouteConfig) SetupHealthRoute() {
	healthRoutes := c.App.Group("/health")
	healthRoutes.Get("/", c.HealthController.GetHealth)
}
