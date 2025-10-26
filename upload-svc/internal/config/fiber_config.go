package config

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:      false,
		AppName:      "upload-svc",
		ErrorHandler: CustomError(),
		BodyLimit:    100 * 1024 * 1024,
		// JSONEncoder:  sonic.Marshal,
		// JSONDecoder:  sonic.Unmarshal,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization, X-Requested-With, Referer, User-Agent",
		ExposeHeaders: "Content-Length",
		MaxAge:        12 * 3600, // 12 hours
	}))

	return app
}

func CustomError() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := http.StatusInternalServerError
		if err, ok := err.(*fiber.Error); ok {
			code = err.Code
		}

		message := &Message{
			Success: false,
			Message: err.Error(),
		}

		return ctx.Status(code).JSON(message)
	}
}

type Message struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
