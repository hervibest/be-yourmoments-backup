package route

import (
	http "github.com/hervibest/be-yourmoments-backup/user-svc/internal/delivery/http/controller"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App              *fiber.App
	AuthController   http.AuthController
	UserController   http.UserController
	ChatController   http.ChatController
	HealthController http.HealthController
	AuthMiddleware   fiber.Handler
}

func (r *RouteConfig) Setup() {
	r.SetupAuthRoute()
	r.SetupUserRoute()
	r.SetupHealthRoute()
}
