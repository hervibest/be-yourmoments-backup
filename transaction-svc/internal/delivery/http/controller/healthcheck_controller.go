package http

import (
	"be-yourmoments/transaction-svc/internal/model"
	"net/http"

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
