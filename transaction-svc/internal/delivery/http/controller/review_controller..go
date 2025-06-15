package http

import (
	"net/http"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type ReviewController interface {
	UserCreateReview(ctx *fiber.Ctx) error
	CreatorGetReview(ctx *fiber.Ctx) error
}

type reviewController struct {
	reviewUseCase   usecase.ReviewUseCase
	customValidator helper.CustomValidator
	logs            *logger.Log
}

func NewReviewController(reviewUseCase usecase.ReviewUseCase,
	customValidator helper.CustomValidator,
	logs *logger.Log) ReviewController {
	return &reviewController{
		reviewUseCase:   reviewUseCase,
		customValidator: customValidator,
		logs:            logs,
	}
}

func (c *reviewController) UserCreateReview(ctx *fiber.Ctx) error {
	request := new(model.CreateReviewRequest)
	if err := ctx.BodyParser(request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	auth := middleware.GetUser(ctx)
	request.UserId = auth.UserId

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.reviewUseCase.Create(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Create review : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.CreatorReviewResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *reviewController) CreatorGetReview(ctx *fiber.Ctx) error {
	request := &model.GetAllReviewRequest{
		Rating: ctx.QueryInt("rating", 0),
		Order:  ctx.Query("order", "DESC"),
		Page:   ctx.QueryInt("page", 1),
		Size:   ctx.QueryInt("size", 10),
	}

	auth := middleware.GetUser(ctx)
	request.CreatorId = auth.CreatorId

	response, pageMetadata, err := c.reviewUseCase.CreatorGetReview(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Get all review error : ", err, c.logs)
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*[]*model.CreatorReviewResponse]{
		Success:      true,
		Data:         response,
		PageMetadata: pageMetadata,
	})
}
