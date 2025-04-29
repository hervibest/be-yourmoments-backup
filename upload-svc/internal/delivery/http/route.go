package http

import (
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/config"

	"github.com/gofiber/fiber/v2"
)

func (c *photoController) PhotoRoute(app *fiber.App, authMiddleware fiber.Handler) {
	api := app.Group(config.EndpointPrefix, authMiddleware)
	api.Post("/single", c.UploadPhoto)
}

func (c *facecamController) FacecamRoute(app *fiber.App, authMiddleware fiber.Handler) {

	// userRoutes := c.App.Group("/api/users", c.AuthMiddleware)
	// userRoutes.Get("/current", c.AuthController.Current)
	// userRoutes.Delete("/logout", c.AuthController.Logout)
	api := app.Group(config.EndpointPrefix, authMiddleware)
	api.Post("/facecam/single", c.UploadFacecam)
}
