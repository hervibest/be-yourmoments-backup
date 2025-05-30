package http

import (
	"net/http"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/usecase"
	"github.com/oklog/ulid/v2"

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
		return helper.ErrUseCaseResponseJSON(ctx, "Create bank : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.BankResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *bankController) FindBankById(ctx *fiber.Ctx) error {
	request := new(model.FindBankByIdRequest)
	request.Id = ctx.Params("bankId")

	if _, err := ulid.Parse(request.Id); err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, "The provided Bank ID is not valid")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.bankUseCase.FindById(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Find bank by id : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.BankResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *bankController) FindAllBank(ctx *fiber.Ctx) error {
	response, err := c.bankUseCase.FindAll(ctx.Context())
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Find all bank : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*[]*model.BankResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *bankController) DeleteBank(ctx *fiber.Ctx) error {
	request := new(model.DeleteBankRequest)
	request.Id = ctx.Params("bankId")
	if _, err := ulid.Parse(request.Id); err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, "The provided Bank ID is not valid")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.bankUseCase.Delete(ctx.Context(), request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Delete bank : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*[]*model.BankResponse]{
		Success: true,
	})
}
