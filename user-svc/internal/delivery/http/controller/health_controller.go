package http

import "github.com/gofiber/fiber/v2"

type HealthController interface {
	GetHealth(ctx *fiber.Ctx) error
}
type healthController struct{}

func NewHealthController() HealthController {
	return &healthController{}
}

func (c *healthController) GetHealth(ctx *fiber.Ctx) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "UP",
		"message": "User service is running",
	})
}
