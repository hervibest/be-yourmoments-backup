package route

import (
	http "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/delivery/http/controller"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                      *fiber.App
	ExploreController        http.ExploreController
	HealthCheckController    http.HealthCheckController
	CheckoutController       http.CheckoutController
	CreatorDiscountControler http.CreatorDiscountController
	PhotoController          http.PhotoController
	AuthMiddleware           fiber.Handler
	CreatorMiddleware        fiber.Handler //TODO doesnt used
}

func (r *RouteConfig) Setup() {
	r.SetupExploreRoute()
	r.SetupHealtCheckRoute()
	r.SetupDiscountRoute()
	r.SetupCheckoutRoute()
	r.SetupPhotoRoute()
}
