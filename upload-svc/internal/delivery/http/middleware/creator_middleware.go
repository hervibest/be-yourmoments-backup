package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/logger"
)

// TODO SHOULD TOKEN VALIDATED ?
func NewCreatorMiddleware(photoAdapter adapter.PhotoAdapter, logs logger.Log) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user := GetUser(ctx)

		creator, err := photoAdapter.GetCreator(ctx.UserContext(), user.UserId)
		if err != nil {
			return helper.ErrUseCaseResponseJSON(ctx, "Authenticate creator error : ", err, logs)
		}

		user.CreatorId = creator.Id

		ctx.Locals("auth", user)
		return ctx.Next()
	}
}
