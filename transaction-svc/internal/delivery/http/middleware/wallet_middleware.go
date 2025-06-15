package middleware

import (
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TODO SHOULD TOKEN VALIDATED ?
func NewWalletMiddleware(walletUseCase usecase.WalletUseCase, tracer trace.Tracer, logs *logger.Log) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logs.Log("wallet middleware called")
		context, span := tracer.Start(ctx.Context(), "GetWallet", oteltrace.WithAttributes())
		defer span.End()

		user := GetUser(ctx)
		request := &model.GetWalletIdRequest{
			CreatorId: user.CreatorId,
		}

		walletId, err := walletUseCase.GetWalletId(context, request)
		if err != nil {
			return helper.ErrUseCaseResponseJSON(ctx, "Get wallet error : ", err, logs)
		}

		user.WalletId = walletId

		ctx.Locals("auth", user)
		return ctx.Next()
	}
}
