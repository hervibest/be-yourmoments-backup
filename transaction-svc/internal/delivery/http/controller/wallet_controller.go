package http

import (
	"net/http"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/usecase"

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
