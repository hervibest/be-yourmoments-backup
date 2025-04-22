package middleware

import (
	"be-yourmoments/user-svc/internal/helper"
	"be-yourmoments/user-svc/internal/helper/logger"
	"be-yourmoments/user-svc/internal/model"
	"be-yourmoments/user-svc/internal/usecase"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// TODO SHOULD TOKEN VALIDATED ?
func NewUserAuth(authUseCase usecase.AuthUseCase, validator helper.CustomValidator, logs *logger.Log) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := strings.TrimPrefix(ctx.Get("Authorization", ""), "Bearer ")
		if token == "" || token == "NOT_FOUND" {
			return fiber.NewError(http.StatusUnauthorized, "Unauthorized access")
		}

		request := &model.VerifyUserRequest{Token: token}

		if validatonErrs := validator.ValidateUseCase(request); validatonErrs != nil {
			return helper.ErrValidationResponseJSON(ctx, validatonErrs)
		}

		auth, err := authUseCase.Verify(ctx.UserContext(), request)
		if err != nil {
			return helper.ErrUseCaseResponseJSON(ctx, err, logs)
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.AuthResponse {
	return ctx.Locals("auth").(*model.AuthResponse)
}
