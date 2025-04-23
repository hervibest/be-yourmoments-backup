package http

import (
	"be-yourmoments/photo-svc/internal/delivery/http/middleware"
	"be-yourmoments/photo-svc/internal/helper"
	"be-yourmoments/photo-svc/internal/helper/logger"
	"be-yourmoments/photo-svc/internal/model"
	"be-yourmoments/photo-svc/internal/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

type CreatorDiscountController interface {
	ActivateDiscount(ctx *fiber.Ctx) error
	CreateDiscount(ctx *fiber.Ctx) error
	DeactivateDiscount(ctx *fiber.Ctx) error
	GetDiscount(ctx *fiber.Ctx) error
}

type creatorDiscountController struct {
	creatorDiscountUseCase usecase.CreatorDiscountUseCase
	customValidator        helper.CustomValidator
	logs                   *logger.Log
}

func NewCreatorDiscountController(creatorDiscountUseCase usecase.CreatorDiscountUseCase, customValidator helper.CustomValidator, logs *logger.Log) CreatorDiscountController {
	return &creatorDiscountController{
		creatorDiscountUseCase: creatorDiscountUseCase,
		customValidator:        customValidator,
		logs:                   logs,
	}
}

func (c *creatorDiscountController) CreateDiscount(ctx *fiber.Ctx) error {
	request := new(model.CreateCreatorDiscountRequest)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	creator := middleware.GetCreator(ctx)
	request.CreatorId = creator.Id

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.creatorDiscountUseCase.CreateDiscount(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Create discount : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.CreatorDiscountResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *creatorDiscountController) ActivateDiscount(ctx *fiber.Ctx) error {
	discountId := ctx.Params("discountId")
	request := &model.ActivateCreatorDiscountRequest{
		Id: discountId,
	}

	creator := middleware.GetCreator(ctx)
	request.CreatorId = creator.Id

	if _, err := ulid.Parse(request.Id); err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, "The provided discount ID is not valid")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.creatorDiscountUseCase.ActivateDiscount(ctx.Context(), request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Activate discount : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *creatorDiscountController) DeactivateDiscount(ctx *fiber.Ctx) error {
	discountId := ctx.Params("discountId")
	request := &model.DeactivateCreatorDiscountRequest{
		Id: discountId,
	}

	creator := middleware.GetCreator(ctx)
	request.CreatorId = creator.Id

	if _, err := ulid.Parse(request.Id); err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, "The provided discount ID is not valid")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.creatorDiscountUseCase.DeactivateDiscount(ctx.Context(), request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Deactivate discount : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *creatorDiscountController) GetDiscount(ctx *fiber.Ctx) error {
	discountId := ctx.Params("discountId")

	request := &model.GetCreatorDiscountRequest{
		Id: discountId,
	}

	creator := middleware.GetCreator(ctx)
	request.CreatorId = creator.Id

	if _, err := ulid.Parse(request.Id); err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, "The provided discount ID is not valid")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.creatorDiscountUseCase.GetDiscount(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Get discount : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.CreatorDiscountResponse]{
		Success: true,
		Data:    response,
	})
}
