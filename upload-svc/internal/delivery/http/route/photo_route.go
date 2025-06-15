package route

import (
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/config"
)

func (r *routeConfig) SetupPhotoRoute() {
	api := r.app.Group(config.EndpointPrefix, r.authMiddleware, r.creatorMiddleware)
	api.Post("/single", r.photoController.UploadPhoto)
	api.Post("/bulk", r.photoController.BulkUploadPhoto)
}
