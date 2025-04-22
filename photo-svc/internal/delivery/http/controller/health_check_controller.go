package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type healthCheckController struct {
}

type HealthCheckController interface {
	HealthCheck(ctx *fiber.Ctx) error
}

func NewHealthCheckController() HealthCheckController {
	return &healthCheckController{}
}

func (c *healthCheckController) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "pong",
	})
}
