package route

import (
	http "github.com/hervibest/be-yourmoments-backup/upload-svc/internal/delivery/http/controller"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig interface {
	Setup()
}

type routeConfig struct {
	app               *fiber.App
	photoController   http.PhotoController
	facecamController http.FacecamController
	authMiddleware    fiber.Handler
}

func NewRouteConfig(app *fiber.App,
	photoController http.PhotoController,
	facecamController http.FacecamController,
	authMiddleware fiber.Handler,
) RouteConfig {
	return &routeConfig{
		app:               app,
		photoController:   photoController,
		facecamController: facecamController,
		authMiddleware:    authMiddleware,
	}
}

func (r *routeConfig) Setup() {
	r.SetupPhotoRoute()
	r.SetupFacecamController()
}
