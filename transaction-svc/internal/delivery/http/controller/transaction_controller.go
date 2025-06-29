package http

import (
	"fmt"
	"net/http"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model/converter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/usecase/contract"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TransactionController interface {
	CreateTransaction(ctx *fiber.Ctx) error
	CreateTransactionV2(ctx *fiber.Ctx) error
	Notify(ctx *fiber.Ctx) error
	GetUserTransactionWithDetail(ctx *fiber.Ctx) error
	GetAllUserTransaction(ctx *fiber.Ctx) error
}

type transactionController struct {
	transactionUseCase contract.TransactionUseCase
	customValidator    helper.CustomValidator
	logs               *logger.Log
}

func NewTransactionController(transactionUseCase contract.TransactionUseCase, customValidator helper.CustomValidator, logs *logger.Log) TransactionController {
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
	request.CreatorId = auth.CreatorId
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if err := helper.MultipleULIDSliceParser(request.PhotoIds); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Invalid photo Ids : ", err, c.logs)
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.transactionUseCase.CreateTransaction(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Create transaction error : ", err, c.logs)
	}

	var httpStatus = http.StatusCreated
	if (response.RedirectURL == "") || (response.SnapToken == "") {
		httpStatus = http.StatusAccepted
	}

	return ctx.Status(httpStatus).JSON(model.WebResponse[*model.CreateTransactionResponse]{
		Success: true,
		Data:    response,
	})
}

// TODO ROBUST VALIDATE FOR EXTERNAL HTTP CALL
func (c *transactionController) Notify(ctx *fiber.Ctx) error {
	webhookRequest := new(model.UpdateTransactionWebhookRequest)

	if err := ctx.BodyParser(webhookRequest); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if _, err := uuid.Parse(webhookRequest.OrderID); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid order id")
	}

	webhookRequest.Body = ctx.Body()

	if validatonErrs := c.customValidator.ValidateUseCase(webhookRequest); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	c.logs.Log(fmt.Sprintf("Received webhook request from midtrans server with fields transaction ID : %s with status : %s ",
		webhookRequest.OrderID, webhookRequest.MidtransTransactionStatus))

	request := converter.WebhookReqToCheckAndUpdate(webhookRequest)
	if err := c.transactionUseCase.CheckAndUpdateTransaction(ctx.Context(), request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Notify webhook : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *transactionController) GetUserTransactionWithDetail(ctx *fiber.Ctx) error {
	request := new(model.GetTransactionWithDetail)
	auth := middleware.GetUser(ctx)

	request.UserID = auth.UserId
	request.TransactionId = ctx.Params("transactionID")

	_, err := uuid.Parse(request.TransactionId)
	if err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, "The provided transaction ID is not valid")
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

func (c *transactionController) CreateTransactionV2(ctx *fiber.Ctx) error {
	request := new(model.CreateTransactionV2Request)
	auth := middleware.GetUser(ctx)

	request.UserId = auth.UserId
	request.CreatorId = auth.CreatorId
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	photoIDs := make([]string, 0, len(request.Items))
	creatorIDs := make([]string, 0, len(request.Items))
	discountIDs := make([]string, 0)
	for _, items := range request.Items {
		photoIDs = append(photoIDs, items.PhotoId)
		creatorIDs = append(photoIDs, items.CreatorId)
		if items.Discount != nil {
			discountIDs = append(discountIDs, items.Discount.DiscountId)
		}
	}

	if err := helper.MultipleULIDSliceParser(photoIDs); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Invalid photo Ids : ", err, c.logs)
	}

	if err := helper.MultipleULIDSliceParser(creatorIDs); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Invalid creator Ids : ", err, c.logs)
	}

	if len(discountIDs) != 0 {
		if err := helper.MultipleULIDSliceParser(discountIDs); err != nil {
			return helper.ErrUseCaseResponseJSON(ctx, "Invalid discount Ids : ", err, c.logs)
		}
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.transactionUseCase.CreateTransactionV2(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Create transaction error : ", err, c.logs)
	}

	var httpStatus = http.StatusCreated
	if (response.RedirectURL == "") || (response.SnapToken == "") {
		httpStatus = http.StatusAccepted
	}

	return ctx.Status(httpStatus).JSON(model.WebResponse[*model.CreateTransactionResponse]{
		Success: true,
		Data:    response,
	})

}
