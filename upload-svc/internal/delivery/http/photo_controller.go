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
	BulkUploadPhoto(ctx *fiber.Ctx) error
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
		return helper.ErrUseCaseResponseJSON(ctx, "Upload photo : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (c *photoController) BulkUploadPhoto(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil || form.File["photo"] == nil {
		return fiber.NewError(http.StatusBadRequest, "No photo files uploaded")
	}

	files := form.File["photo"]

	request := new(model.CreatePhotoRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid body json")
	}

	auth := middleware.GetUser(ctx)

	request.UserId = auth.UserId
	request.CreatorId = auth.CreatorId

	priceStr := strconv.Itoa(request.Price)
	request.PriceStr = priceStr

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	err = c.photoUsecase.BulkUploadPhoto(ctx.UserContext(), files, request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Upload bulk photo photo : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}
