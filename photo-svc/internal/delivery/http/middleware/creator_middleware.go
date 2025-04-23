package middleware

import (
	"be-yourmoments/photo-svc/internal/helper"
	"be-yourmoments/photo-svc/internal/helper/logger"
	"be-yourmoments/photo-svc/internal/model"
	"be-yourmoments/photo-svc/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TODO SHOULD TOKEN VALIDATED ?
func NewCreatorMiddleware(creatorUseCase usecase.CreatorUseCase, tracer trace.Tracer, logs *logger.Log) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		context, span := tracer.Start(ctx.Context(), "creatorMiddleware", oteltrace.WithAttributes())
		defer span.End()

		auth := GetUser(ctx)
		request := &model.GetCreatorRequest{
			UserId: auth.UserId,
		}

		creatorResponse, err := creatorUseCase.GetCreator(context, request)
		if err != nil {
			return helper.ErrUseCaseResponseJSON(ctx, "Get creator error : ", err, logs)
		}

		creator := &model.CreatorResponse{
			Id: creatorResponse.Id,
		}

		ctx.Locals("creator", creator)
		return ctx.Next()
	}
}

func GetCreator(ctx *fiber.Ctx) *model.CreatorResponse {
	return ctx.Locals("creator").(*model.CreatorResponse)
}
