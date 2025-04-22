package http

import (
	"be-yourmoments/upload-svc/internal/delivery/http/middleware"
	"be-yourmoments/upload-svc/internal/helper"
	"be-yourmoments/upload-svc/internal/helper/logger"
	"be-yourmoments/upload-svc/internal/model"
	"be-yourmoments/upload-svc/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PhotoController interface {
	UploadPhoto(ctx *fiber.Ctx) error
	PhotoRoute(app *fiber.App, authMiddleware fiber.Handler)
}

type photoController struct {
	photoUsecase    usecase.PhotoUsecase
	logs            *logger.Log
	customValidator helper.CustomValidator
}

func NewPhotoController(photoUsecase usecase.PhotoUsecase, logs *logger.Log, customValidator helper.CustomValidator) PhotoController {
	return &photoController{
		photoUsecase:    photoUsecase,
		logs:            logs,
		customValidator: customValidator,
	}
}

// TODO CUSTOM VALIDATOR
func (c *photoController) UploadPhoto(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("photo")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid photo")
	}

	request := new(model.CreatePhotoRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid body json")
	}

	auth := middleware.GetUser(ctx)
	request.UserId = auth.UserId
	priceStr := strconv.Itoa(request.Price)
	request.PriceStr = priceStr

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	err = c.photoUsecase.UploadPhoto(ctx.UserContext(), file, request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}
