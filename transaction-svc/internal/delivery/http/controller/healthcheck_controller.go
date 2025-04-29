package http

import (
	"net/http"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"

	"github.com/gofiber/fiber/v2"
)

type HealthCheckController interface {
	HealthCheck(ctx *fiber.Ctx) error
}

type healthCheckController struct {
}

func NewHealthCheckController() HealthCheckController {
	return &healthCheckController{}
}

func (c *healthCheckController) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
		Message: "pong",
	})
}
