package middleware

import (
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TODO SHOULD TOKEN VALIDATED ?
func NewCreatorMiddleware(photoAdapter adapter.PhotoAdapter, tracer trace.Tracer, logs *logger.Log) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logs.Log("creator middleware called")
		context, span := tracer.Start(ctx.Context(), "authenticateCreator", oteltrace.WithAttributes())
		defer span.End()

		user := GetUser(ctx)

		creator, err := photoAdapter.GetCreator(context, user.UserId)
		if err != nil {
			return helper.ErrUseCaseResponseJSON(ctx, "Authenticate user error : ", err, logs)
		}

		user.CreatorId = creator.Id

		ctx.Locals("auth", user)
		return ctx.Next()
	}
}
