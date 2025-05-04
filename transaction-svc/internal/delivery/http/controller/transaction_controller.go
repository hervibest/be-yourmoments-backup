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
	"github.com/google/uuid"
)

type TransactionController interface {
	CreateTransaction(ctx *fiber.Ctx) error
	Notify(ctx *fiber.Ctx) error
	GetUserTransactionWithDetail(ctx *fiber.Ctx) error
	GetAllUserTransaction(ctx *fiber.Ctx) error
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

// TODO ROBUST VALIDATE FOR EXTERNAL HTTP CALL
func (c *transactionController) GetUserTransactionWithDetail(ctx *fiber.Ctx) error {
	request := new(model.GetTransactionWithDetail)
	auth := middleware.GetUser(ctx)

	request.UserID = auth.UserId
	c.logs.Log(request.UserID)
	request.TransactionId = ctx.Params("transactionID")

	if _, err := uuid.Parse(request.TransactionId); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid transaction id")
	}

	if _, err := ulid.Parse(request.UserID); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid user id")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.transactionUseCase.UserGetWithDetail(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "User Get With Detail : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.TransactionWithDetail]{
		Success: true,
		Data:    response,
	})
}

func (c *transactionController) GetAllUserTransaction(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.GetAllUsertTransaction{
		UserId: auth.UserId,
		Order:  ctx.Query("order", "DESC"),
		Page:   ctx.QueryInt("page", 1),
		Size:   ctx.QueryInt("size", 10),
	}

	response, pageMetadata, err := c.transactionUseCase.GetAllUserTransaction(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Get all user transaction error : ", err, c.logs)
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*[]*model.UserTransaction]{
		Success:      true,
		Data:         response,
		PageMetadata: pageMetadata,
	})
}
