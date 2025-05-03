package http

import (
	"context"
	"log"
	"net/http"
	"path/filepath"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/usecase"
	"github.com/oklog/ulid/v2"

	"github.com/gofiber/fiber/v2"
)

type PhotoController interface {
	GetBulkPhotoDetail(ctx *fiber.Ctx) error
	GetPhotoFile(ctx *fiber.Ctx) error
}
type photoController struct {
	photoUseCase    usecase.PhotoUseCase
	customValidator helper.CustomValidator
	logs            *logger.Log
}

func NewPhotoController(photoUseCase usecase.PhotoUseCase, customValidator helper.CustomValidator, logs *logger.Log) PhotoController {
	return &photoController{photoUseCase: photoUseCase, customValidator: customValidator, logs: logs}
}

func (c *photoController) GetBulkPhotoDetail(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.GetBulkPhotoDetailRequest)
	request.CreatorId = auth.CreatorId
	request.BulkPhotoId = ctx.Params("bulkPhotoId")

	if _, err := ulid.Parse(request.BulkPhotoId); err != nil {
		return fiber.NewError(http.StatusBadRequest, "The provided photo ID is not valid")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.photoUseCase.GetBulkPhotoDetail(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Get all explore similar error : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.GetBulkPhotoDetailResponse]{
		Success: true,
		Data:    response,
	})
}

// TODO security, validation and parsing logic (exposed for CDN)
func (c *photoController) GetPhotoFile(ctx *fiber.Ctx) error {
	log.Print("accessed")

	fileKey := ctx.Query("fileKey")

	object, err := c.photoUseCase.GetPhotoFile(context.Background(), fileKey)
	if err != nil {
		c.logs.Log(err)
		return ctx.Status(500).SendString("Error getting object")
	}

	// Optional: deteksi tipe file dari extensi
	contentType := "image/jpeg"
	if ext := filepath.Ext(fileKey); ext == ".png" {
		contentType = "image/png"
	}

	ctx.Set("Content-Type", contentType)
	ctx.Set("Cache-Control", "public, max-age=3600") //caching ??
	return ctx.SendStream(object)
}
