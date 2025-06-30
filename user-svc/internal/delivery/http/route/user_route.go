package route

func (c *RouteConfig) SetupUserRoute() {
	userRoutesV2 := c.App.Group("/api/v2/users", c.AuthMiddleware)
	userRoutesV2.Get("/profile", c.UserController.GetUserProfileV2)
	userRoutesV2.Patch("/profile", c.UserController.UpdateUserProfileImageV2)
	userRoutesV2.Patch("/profile/cover", c.UserController.UpdateUserCoverImageV2)

	userRoutes := c.App.Group("/api/users", c.AuthMiddleware)
	userRoutes.Get("/current", c.AuthController.Current)
	userRoutes.Delete("/logout", c.AuthController.Logout)
	userRoutes.Get("/dm", c.UserController.GetAllPublicUserChat)

	userRoutes.Post("/room", c.ChatController.GetOrCreateRoom)
	userRoutes.Get("/token/:uid", c.ChatController.GetCustomToken)
	userRoutes.Post("/send-message", c.ChatController.SendMessage)

	userRoutes.Put("/similarity", c.UserController.UpdateUserSimilarity)
}
