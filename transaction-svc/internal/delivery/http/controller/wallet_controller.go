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

type WalletController interface {
	GetWallet(ctx *fiber.Ctx) error
}
type walletController struct {
	walletUseCase usecase.WalletUsecase
	logs          *logger.Log
}

func NewWalletController(walletUseCase usecase.WalletUsecase, logs *logger.Log) WalletController {
	return &walletController{walletUseCase: walletUseCase, logs: logs}
}

func (c *walletController) GetWallet(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetWalletRequest{
		CreatorId: auth.CreatorId,
	}

	wallet, err := c.walletUseCase.GetWallet(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Create withdrawal : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.WalletResponse]{
		Success: true,
		Data:    wallet,
	})
}
