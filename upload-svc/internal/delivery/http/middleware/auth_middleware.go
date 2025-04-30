package middleware

import (
	"strings"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/model"

	"github.com/gofiber/fiber/v2"
)

// TODO SHOULD TOKEN VALIDATED ?
func NewUserAuth(userAdapter adapter.UserAdapter, logs logger.Log) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := strings.TrimPrefix(ctx.Get("Authorization", ""), "Bearer ")
		if token == "" || token == "NOT_FOUND" {
			return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized access")
		}

		authResponse, err := userAdapter.AuthenticateUser(ctx.UserContext(), token)
		if err != nil {
			return helper.ErrUseCaseResponseJSON(ctx, "Authenticate user : ", err, logs)
		}

		auth := &model.AuthResponse{
			UserId:      authResponse.GetUser().GetUserId(),
			Username:    authResponse.GetUser().GetUsername(),
			Email:       authResponse.GetUser().GetEmail(),
			PhoneNumber: authResponse.GetUser().GetPhoneNumber(),
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
