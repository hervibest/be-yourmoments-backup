package http

import (
	"be-yourmoments/transaction-svc/internal/helper"
	"be-yourmoments/transaction-svc/internal/helper/logger"
	"be-yourmoments/transaction-svc/internal/model"
	"be-yourmoments/transaction-svc/internal/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type BankController interface {
	CreateBank(ctx *fiber.Ctx) error
	DeleteBank(ctx *fiber.Ctx) error
	FindAllBank(ctx *fiber.Ctx) error
	FindBankById(ctx *fiber.Ctx) error
}

type bankController struct {
	bankUseCase     usecase.BankUseCase
	customValidator helper.CustomValidator
	logs            *logger.Log
}

func NewBankController(bankUseCase usecase.BankUseCase,
	customValidator helper.CustomValidator,
	logs *logger.Log) BankController {
	return &bankController{
		bankUseCase:     bankUseCase,
		customValidator: customValidator,
		logs:            logs,
	}
}

func (c *bankController) CreateBank(ctx *fiber.Ctx) error {
	request := new(model.CreateBankRequest)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.bankUseCase.Create(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.BankResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *bankController) FindBankById(ctx *fiber.Ctx) error {
	request := new(model.FindBankByIdRequest)
	request.Id = ctx.Params("bankId")

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.bankUseCase.FindById(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.BankResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *bankController) FindAllBank(ctx *fiber.Ctx) error {
	response, err := c.bankUseCase.FindAll(ctx.Context())
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*[]*model.BankResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *bankController) DeleteBank(ctx *fiber.Ctx) error {
	request := new(model.DeleteBankRequest)
	request.Id = ctx.Params("bankId")
	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.bankUseCase.Delete(ctx.Context(), request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*[]*model.BankResponse]{
		Success: true,
	})
}
