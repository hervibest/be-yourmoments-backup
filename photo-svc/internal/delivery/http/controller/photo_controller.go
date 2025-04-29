package http

import (
	"net/http"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type PhotoController interface {
	GetBulkPhotoDetail(ctx *fiber.Ctx) error
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
