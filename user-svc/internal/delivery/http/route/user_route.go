package route

func (c *RouteConfig) SetupUserRoute() {
	userRoutes := c.App.Group("/api/users", c.AuthMiddleware)
	userRoutes.Get("/current", c.AuthController.Current)
	userRoutes.Delete("/logout", c.AuthController.Logout)
	userRoutes.Get("/profile", c.UserController.GetUserProfile)

	userRoutes.Put("/profile", c.UserController.UpdateUserProfile)
	userRoutes.Post("/profile/upload-profile-image", c.UserController.UploadUserProfileImage)
	userRoutes.Post("/profile/upload-profile-cover", c.UserController.UploadUserCoverImage)
	userRoutes.Get("/dm", c.UserController.GetAllPublicUserChat)

	userRoutes.Post("/room", c.ChatController.GetOrCreateRoom)
	userRoutes.Get("/token/:uid", c.ChatController.GetCustomToken)
	userRoutes.Post("/send-message", c.ChatController.SendMessage)

	userRoutes.Put("/similarity", c.UserController.UpdateUserSimilarity)
}
