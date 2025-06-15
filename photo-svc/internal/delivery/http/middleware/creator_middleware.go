package middleware

import (
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func NewCreatorMiddleware(creatorUseCase usecase.CreatorUseCase, tracer trace.Tracer, logs *logger.Log) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		context, span := tracer.Start(ctx.Context(), "creatorMiddleware", oteltrace.WithAttributes())
		defer span.End()

		auth := GetUser(ctx)
		request := &model.GetCreatorIdRequest{
			UserId: auth.UserId,
		}

		creatorId, err := creatorUseCase.GetCreatorId(context, request)
		if err != nil {
			return helper.ErrUseCaseResponseJSON(ctx, "Get creator error : ", err, logs)
		}

		auth.CreatorId = creatorId

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}
