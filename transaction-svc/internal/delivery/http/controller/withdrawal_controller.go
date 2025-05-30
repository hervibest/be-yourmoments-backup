package http

import (
	"net/http"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/usecase"
	"github.com/oklog/ulid/v2"

	"github.com/gofiber/fiber/v2"
)

type WithdrawalController interface {
	CreateWithdrawal(ctx *fiber.Ctx) error
	// DeleteWithdrawal(ctx *fiber.Ctx) error
	FindAllWithdrawal(ctx *fiber.Ctx) error
	FindWithdrawalById(ctx *fiber.Ctx) error
}

type withdrawalController struct {
	withdrawalUseCase usecase.WithdrawalUseCase
	customValidator   helper.CustomValidator
	logs              *logger.Log
}

func NewWithdrawalController(withdrawalUseCase usecase.WithdrawalUseCase,
	customValidator helper.CustomValidator,
	logs *logger.Log) WithdrawalController {
	return &withdrawalController{
		withdrawalUseCase: withdrawalUseCase,
		customValidator:   customValidator,
		logs:              logs,
	}
}

func (c *withdrawalController) CreateWithdrawal(ctx *fiber.Ctx) error {
	request := new(model.CreateWithdrawalRequest)
	auth := middleware.GetUser(ctx)
	request.WalletId = auth.WalletId

	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	ulidErrMaps := map[string]string{
		request.WalletId:     "The provided Wallet ID is not valid",
		request.BankWalletId: "The provided Bank Wallet ID is not valid",
	}

	if err := helper.MultipleULIDParser(ulidErrMaps); err != nil {
		return err
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.withdrawalUseCase.Create(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Create withdrawal : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.WithdrawalResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *withdrawalController) FindWithdrawalById(ctx *fiber.Ctx) error {
	request := new(model.FindWithdrawalById)
	request.Id = ctx.Params("withdrawalId")
	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if _, err := ulid.Parse(request.Id); err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, "The provided Withdrawal ID is not valid")
	}

	response, err := c.withdrawalUseCase.FindById(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Find withdrawal by id : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.WithdrawalResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *withdrawalController) FindAllWithdrawal(ctx *fiber.Ctx) error {
	response, err := c.withdrawalUseCase.FindAll(ctx.Context())
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Find all withdrawal : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*[]*model.WithdrawalResponse]{
		Success: true,
		Data:    response,
	})
}

// func (c *withdrawalController) DeleteWithdrawal(ctx *fiber.Ctx) error {
// 	request := new(model.DeleteWithdrawalRequest)
// 	if err := helper.StrictBodyParser(ctx, request); err != nil {
// 		return helper.ErrBodyParserResponseJSON(ctx, err)
// 	}

// 	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
// 		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
// 	}

// 	if err := c.withdrawalUseCase.Delete(ctx.Context(), request); err != nil {
// 		return helper.ErrUseCaseResponseJSON(ctx, err, c.logs)
// 	}

// 	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*[]*model.WithdrawalResponse]{
// 		Success: true,
// 	})
// }
