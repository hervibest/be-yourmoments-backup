package http

import (
	"net/http"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TransactionController interface {
	CreateTransaction(ctx *fiber.Ctx) error
	Notify(ctx *fiber.Ctx) error
}

type transactionController struct {
	transactionUseCase usecase.TransactionUseCase
	customValidator    helper.CustomValidator
	logs               *logger.Log
}

func NewTransactionController(transactionUseCase usecase.TransactionUseCase, customValidator helper.CustomValidator, logs *logger.Log) TransactionController {
	return &transactionController{
		transactionUseCase: transactionUseCase,
		customValidator:    customValidator,
		logs:               logs,
	}
}

func (c *transactionController) CreateTransaction(ctx *fiber.Ctx) error {
	request := new(model.CreateTransactionRequest)
	auth := middleware.GetUser(ctx)

	request.UserId = auth.UserId
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.transactionUseCase.CreateTransaction(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Create transaction error : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.CreateTransactionResponse]{
		Success: true,
		Data:    response,
	})
}

// TODO ROBUST VALIDATE FOR EXTERNAL HTTP CALL
func (c *transactionController) Notify(ctx *fiber.Ctx) error {
	request := new(model.UpdateTransactionWebhookRequest)

	if err := ctx.BodyParser(request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if _, err := uuid.Parse(request.OrderID); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid order id")
	}

	request.Body = ctx.Body()

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.transactionUseCase.UpdateTransactionWebhook(ctx.Context(), request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Notify webhook : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[any]{
		Success: true,
	})
}
