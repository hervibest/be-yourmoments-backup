package route

import (
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/config"
)

func (r *routeConfig) SetupFacecamController() {
	// userRoutes := c.App.Group("/api/users", c.AuthMiddleware)
	// userRoutes.Get("/current", c.AuthController.Current)
	// userRoutes.Delete("/logout", c.AuthController.Logout)
	api := r.app.Group(config.EndpointPrefix, r.authMiddleware, r.creatorMiddleware)
	api.Post("/facecam/single", r.facecamController.UploadFacecam)
}
