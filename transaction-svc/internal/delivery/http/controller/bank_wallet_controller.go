package http

import (
	"net/http"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type BankWalletController interface {
	CreateBankWallet(ctx *fiber.Ctx) error
	DeleteBankWallet(ctx *fiber.Ctx) error
	FindAllBankWallet(ctx *fiber.Ctx) error
}

type bankWalletController struct {
	bankWalletUseCase usecase.BankWalletUseCase
	customValidator   helper.CustomValidator
	logs              *logger.Log
}

func NewBankWalletController(bankWalletUseCase usecase.BankWalletUseCase,
	customValidator helper.CustomValidator,
	logs *logger.Log) BankWalletController {
	return &bankWalletController{
		bankWalletUseCase: bankWalletUseCase,
		customValidator:   customValidator,
		logs:              logs,
	}
}

func (c *bankWalletController) CreateBankWallet(ctx *fiber.Ctx) error {
	request := new(model.CreateBankWalletRequest)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.bankWalletUseCase.Create(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Create bank wallet : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.BankWalletResponse]{
		Success: true,
		Data:    response,
	})
}

// func (c *bankWalletController) FindBankById(ctx *fiber.Ctx) error {
// 	request := new(model.FindByIdRequest)
// 	if err := helper.StrictBodyParser(ctx, request); err != nil {
// 		return helper.ErrBodyParserResponseJSON(ctx, err)
// 	}

// 	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
// 		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
// 	}

// 	response, err := c.bankWalletUseCase.FindById(ctx.Context(), request)
// 	if err != nil {
// 		return helper.ErrUseCaseResponseJSON(ctx, err, c.logs)
// 	}

// 	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.BankWalletResponse]{
// 		Success: true,
// 		Data:    response,
// 	})
// }

func (c *bankWalletController) FindAllBankWallet(ctx *fiber.Ctx) error {
	response, err := c.bankWalletUseCase.FindAll(ctx.Context())
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Find all bank wallet : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*[]*model.BankWalletResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *bankWalletController) DeleteBankWallet(ctx *fiber.Ctx) error {
	request := new(model.DeleteBankWalletRequest)
	request.Id = ctx.Params("bankWalletId")

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.bankWalletUseCase.Delete(ctx.Context(), request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Delete bank wallet : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*[]*model.BankWalletResponse]{
		Success: true,
	})
}
