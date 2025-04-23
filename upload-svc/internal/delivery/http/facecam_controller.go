package http

import (
	"be-yourmoments/upload-svc/internal/delivery/http/middleware"
	"be-yourmoments/upload-svc/internal/helper"
	"be-yourmoments/upload-svc/internal/helper/logger"
	"be-yourmoments/upload-svc/internal/usecase"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type FacecamController interface {
	UploadFacecam(ctx *fiber.Ctx) error
	FacecamRoute(app *fiber.App, authMiddleware fiber.Handler)
}

type facecamController struct {
	facecamUseCase usecase.FacecamUseCase
	logs           *logger.Log
}

func NewFacecamController(facecamUseCase usecase.FacecamUseCase, logs *logger.Log) FacecamController {
	return &facecamController{
		facecamUseCase: facecamUseCase,
		logs:           logs,
	}
}

func (c *facecamController) UploadFacecam(ctx *fiber.Ctx) error {
	log.Println("Upload facecam via http")
	file, err := ctx.FormFile("facecam")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid facecam")
	}

	auth := middleware.GetUser(ctx)

	err = c.facecamUseCase.UploadFacecam(ctx.UserContext(), file, auth.UserId)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Upload facecam : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
	})

}
