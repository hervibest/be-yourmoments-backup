package http

import (
	"be-yourmoments/transaction-svc/internal/delivery/http/middleware"
	"be-yourmoments/transaction-svc/internal/helper"
	"be-yourmoments/transaction-svc/internal/helper/logger"
	"be-yourmoments/transaction-svc/internal/model"
	"be-yourmoments/transaction-svc/internal/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type TransactionWalletController interface {
	GetAll(ctx *fiber.Ctx) error
}

type transactionWalletController struct {
	transactionWalletUC usecase.TransactionWalletUseCase
	customValidator     helper.CustomValidator
	logs                *logger.Log
}

func NewTransactionWalletController(transactionWalletUC usecase.TransactionWalletUseCase, customValidator helper.CustomValidator, logs *logger.Log) TransactionWalletController {
	return &transactionWalletController{transactionWalletUC: transactionWalletUC, customValidator: customValidator, logs: logs}
}

func (c *transactionWalletController) GetAll(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.GetAllTransactionWallet{
		WalletId: auth.WalletId,
		Max:      ctx.Query("max", ""),
		Min:      ctx.Query("min", ""),
		Order:    ctx.Query("order", "DESC"),
		Page:     ctx.QueryInt("page", 1),
		Size:     ctx.QueryInt("size", 10),
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, pageMetadata, err := c.transactionWalletUC.GetAll(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Get all wallets error : ", err, c.logs)
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*[]*model.TransactionWalletResponse]{
		Success:      true,
		Data:         response,
		PageMetadata: pageMetadata,
	})

}
