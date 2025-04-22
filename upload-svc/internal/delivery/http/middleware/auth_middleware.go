package middleware

import (
	"be-yourmoments/upload-svc/internal/adapter"
	"be-yourmoments/upload-svc/internal/helper"
	"be-yourmoments/upload-svc/internal/helper/logger"
	"be-yourmoments/upload-svc/internal/model"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// TODO SHOULD TOKEN VALIDATED ?
func NewUserAuth(userAdapter adapter.UserAdapter, logs *logger.Log) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := strings.TrimPrefix(ctx.Get("Authorization", ""), "Bearer ")
		if token == "" || token == "NOT_FOUND" {
			return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized access")
		}

		authResponse, err := userAdapter.AuthenticateUser(ctx.UserContext(), token)
		if err != nil {
			return helper.ErrUseCaseResponseJSON(ctx, err, logs)
		}

		auth := &model.AuthResponse{
			UserId:      authResponse.User.GetUserId(),
			Username:    authResponse.User.GetUsername(),
			Email:       authResponse.User.GetEmail(),
			PhoneNumber: authResponse.User.GetPhoneNumber(),
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.AuthResponse {
	return ctx.Locals("auth").(*model.AuthResponse)
}
