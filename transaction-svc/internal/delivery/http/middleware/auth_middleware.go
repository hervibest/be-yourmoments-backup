package middleware

import (
	"strings"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TODO SHOULD TOKEN VALIDATED ?
func NewUserAuth(userAdapter adapter.UserAdapter, tracer trace.Tracer, logs *logger.Log) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		context, span := tracer.Start(ctx.Context(), "authenticateUser", oteltrace.WithAttributes())
		defer span.End()
		token := strings.TrimPrefix(ctx.Get("Authorization", ""), "Bearer ")
		if token == "" || token == "NOT_FOUND" {
			return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized access")
		}

		authResponse, err := userAdapter.AuthenticateUser(context, token)
		if err != nil {
			return helper.ErrUseCaseResponseJSON(ctx, "Authenticate user error : ", err, logs)
		}

		auth := &model.AuthResponse{
			UserId:      authResponse.GetUser().GetUserId(),
			Username:    authResponse.GetUser().GetUsername(),
			Email:       authResponse.GetUser().GetEmail(),
			PhoneNumber: authResponse.GetUser().GetPhoneNumber(),
			Similarity:  authResponse.GetUser().GetSimilarity(),
			CreatorId:   authResponse.GetUser().GetCreatorId(),
			WalletId:    authResponse.GetUser().GetWalletId(),
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.AuthResponse {
	return ctx.Locals("auth").(*model.AuthResponse)
}
