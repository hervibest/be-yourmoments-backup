package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type PhotoController interface {
	UploadPhoto(ctx *fiber.Ctx) error
	BulkUploadPhoto(ctx *fiber.Ctx) error
}

type photoController struct {
	photoUsecase    usecase.PhotoUsecase
	logs            logger.Log
	customValidator helper.CustomValidator
}

func NewPhotoController(photoUsecase usecase.PhotoUsecase, logs logger.Log, customValidator helper.CustomValidator) PhotoController {
	return &photoController{
		photoUsecase:    photoUsecase,
		logs:            logs,
		customValidator: customValidator,
	}
}

// -- IMPLEMENT COMPRESS PHOTO (DEFAULT)
func (c *photoController) UploadPhoto(ctx *fiber.Ctx) error {
	startController := time.Now()
	file, err := ctx.FormFile("photo")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "No photo file found in the request.")
	}

	c.logs.Log(fmt.Sprintf("⏱️ File check form file: %v", time.Since(startController)))

	request := new(model.CreatePhotoRequest)
	if err := ctx.BodyParser(request); err != nil {
		return helper.ErrFormParserResponseJSON(ctx, "Failed to parse form fields. Please ensure all fields are sent in correct format.", err, c.logs)
	}

	auth := middleware.GetUser(ctx)
	request.UserId = auth.UserId
	request.CreatorId = auth.CreatorId
	priceStr := strconv.Itoa(request.Price)
	request.PriceStr = priceStr

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}
	c.logs.Log(fmt.Sprintf("⏱️ File check validate file: %v", time.Since(startController)))

	err = c.photoUsecase.UploadPhoto(ctx.UserContext(), file, request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Upload photo : ", err, c.logs)
	}

	c.logs.Log(fmt.Sprintf("⏱️ File check total file: %v", time.Since(startController)))

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

// -- NOT IMPLEMENTING COMPRESS PHOTO (DEFAULT)
/*  CONCERN
1. User bisa saja berhasil upload semua foto tetapi kalau semisal ada validasi bisnis logic (metadata dan harga) gagal apa yang terjadi
2. Perlukah untuk menyimpan resumable upload ? Mengingat bisa sangat bulk alias besar
*/
func (c *photoController) BulkUploadPhoto(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil || form.File["photo"] == nil {
		return fiber.NewError(http.StatusBadRequest, "No photo files found in the request. Make sure to include at least one photo.")
	}

	files := form.File["photo"]

	request := new(model.CreatePhotoRequest)
	if err := ctx.BodyParser(request); err != nil {
		return helper.ErrFormParserResponseJSON(ctx, "Failed to parse form fields. Please ensure all fields are sent in correct format.", err, c.logs)
	}

	auth := middleware.GetUser(ctx)

	request.UserId = auth.UserId
	request.CreatorId = auth.CreatorId

	priceStr := strconv.Itoa(request.Price)
	request.PriceStr = priceStr

	const maxFileSize = 1 * 1024 * 1024 // 1MB
	for _, file := range files {
		if file.Size > maxFileSize {
			return fiber.NewError(fiber.StatusRequestEntityTooLarge, "File size exceeds the 1MB limit")
		}
	}

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
