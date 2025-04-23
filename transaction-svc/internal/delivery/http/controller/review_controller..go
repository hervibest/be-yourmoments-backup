package http

import (
	"be-yourmoments/transaction-svc/internal/helper"
	"be-yourmoments/transaction-svc/internal/helper/logger"
	"be-yourmoments/transaction-svc/internal/model"
	"be-yourmoments/transaction-svc/internal/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ReviewController interface {
	CreateReview(ctx *fiber.Ctx) error
	GetAllReview(ctx *fiber.Ctx) error
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

func (c *reviewController) CreateReview(ctx *fiber.Ctx) error {
	request := new(model.CreateReviewRequest)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

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

func (c *reviewController) GetAllReview(ctx *fiber.Ctx) error {
	request := &model.GetAllReviewRequest{
		Star:  ctx.QueryInt("username", 0),
		Order: ctx.Query("order", "DESC"),
		Page:  ctx.QueryInt("page", 1),
		Size:  ctx.QueryInt("size", 10),
	}

	response, pageMetadata, err := c.reviewUseCase.GetCreatorReview(ctx.Context(), request)
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
