package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/enum/message"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	CreateDeviceToken(ctx *fiber.Ctx) error
	Current(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	RegisterByEmail(ctx *fiber.Ctx) error
	RegisterOrLoginByGoogle(ctx *fiber.Ctx) error
	RegisterByPhoneNumber(ctx *fiber.Ctx) error
	RequestAccessToken(ctx *fiber.Ctx) error
	RequestResetPassword(ctx *fiber.Ctx) error
	ResendEmailVerification(ctx *fiber.Ctx) error
	ResetPassword(ctx *fiber.Ctx) error
	ValidateResetPassword(ctx *fiber.Ctx) error
	VerifyEmail(ctx *fiber.Ctx) error
}

type authController struct {
	authUseCase     usecase.AuthUseCase
	customValidator helper.CustomValidator
	logs            logger.Log
}

func NewAuthController(authUseCase usecase.AuthUseCase, customValidator helper.CustomValidator, logs logger.Log) AuthController {
	return &authController{
		authUseCase:     authUseCase,
		customValidator: customValidator,
		logs:            logs,
	}
}

func (c *authController) RegisterByPhoneNumber(ctx *fiber.Ctx) error {
	request := new(model.RegisterByPhoneRequest)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	parsedDate, err := time.Parse("2006-01-02", request.BirthDateStr)
	if err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, errors.New(message.InvalidBirthDate))
	}

	request.BirthDate = &parsedDate

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.authUseCase.RegisterByPhoneNumber(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Register by phone number error : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.UserResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *authController) RegisterOrLoginByGoogle(ctx *fiber.Ctx) error {
	request := new(model.RegisterByGoogleRequest)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, token, err := c.authUseCase.RegisterOrLoginByGoogle(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Register by google sign in error : ", err, c.logs)
	}

	responses := map[string]interface{}{
		"user":  response,
		"token": token,
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
		Data:    responses,
	})
}

func (c *authController) RegisterByEmail(ctx *fiber.Ctx) error {
	request := new(model.RegisterByEmailRequest)
	if err := ctx.BodyParser(request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	parsedDate, err := time.Parse("2006-01-02", request.BirthDateStr)
	if err != nil {
		return helper.ErrBodyResponseJSON(ctx, message.InvalidBirthDate)
	}

	request.BirthDate = &parsedDate
	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.authUseCase.RegisterByEmail(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Register by email error : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.UserResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *authController) ResendEmailVerification(ctx *fiber.Ctx) error {
	request := new(model.ResendEmailUserRequest)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.authUseCase.ResendEmailVerification(ctx.Context(), request.Email); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Resend email verification error : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *authController) VerifyEmail(ctx *fiber.Ctx) error {
	request := new(model.VerifyEmailUserRequest)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	request.Token = ctx.Params("token")
	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.authUseCase.VerifyEmail(ctx.Context(), request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Verify email : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *authController) RequestResetPassword(ctx *fiber.Ctx) error {
	request := new(model.SendResetPasswordRequest)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.authUseCase.RequestResetPassword(ctx.Context(), request.Email); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Request reset password error : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})

}

func (c *authController) ValidateResetPassword(ctx *fiber.Ctx) error {
	request := new(model.ValidateResetTokenRequest)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	valid, err := c.authUseCase.ValidateResetPassword(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Validate reset password error : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
		Data: map[string]interface{}{
			"valid": valid,
		},
	})
}

func (c *authController) ResetPassword(ctx *fiber.Ctx) error {
	request := new(model.ResetPasswordUserRequest)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	request.Token = ctx.Params("token")
	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.authUseCase.ResetPassword(ctx.Context(), request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Reset password error : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *authController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	userResponse, tokenResponse, err := c.authUseCase.Login(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Login error : ", err, c.logs)
	}

	response := map[string]interface{}{
		"user":  userResponse,
		"token": tokenResponse,
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
		Data:    response,
	})
}

func (c *authController) CreateDeviceToken(ctx *fiber.Ctx) error {
	request := new(model.DeviceRequest)
	if err := ctx.BodyParser(request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	auth := middleware.GetUser(ctx)
	request.UserId = auth.UserId

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.authUseCase.CreateDeviceToken(ctx.Context(), request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Create device token error : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *authController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	userResponse, err := c.authUseCase.Current(ctx.Context(), auth.Email)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Current error : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.UserResponse]{
		Success: true,
		Data:    userResponse,
	})
}

func (c *authController) RequestAccessToken(ctx *fiber.Ctx) error {
	request := new(model.AccessTokenRequest)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	userResponse, tokenResponse, err := c.authUseCase.AccessTokenRequest(ctx.Context(), request.RefreshToken)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Request access token error : ", err, c.logs)
	}

	responses := map[string]interface{}{
		"user":  userResponse,
		"token": tokenResponse,
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
		Data:    responses,
	})
}

func (c *authController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.LogoutUserRequest)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	request.UserId = auth.UserId
	request.AccessToken = auth.Token
	request.ExpiresAt = auth.ExpiresAt

	valid, err := c.authUseCase.Logout(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Logout error : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
		Data: map[string]interface{}{
			"valid": valid,
		},
	})
}
