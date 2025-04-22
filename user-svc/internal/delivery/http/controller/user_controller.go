package http

import (
	"be-yourmoments/user-svc/internal/delivery/http/middleware"
	"be-yourmoments/user-svc/internal/helper"
	"be-yourmoments/user-svc/internal/helper/logger"
	"be-yourmoments/user-svc/internal/model"
	"be-yourmoments/user-svc/internal/usecase"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	GetUserProfile(ctx *fiber.Ctx) error
	UpdateUserProfile(ctx *fiber.Ctx) error
	UpdateUserProfileImage(ctx *fiber.Ctx) error
	UpdateUserCoverImage(ctx *fiber.Ctx) error
	GetAllPublicUserChat(ctx *fiber.Ctx) error
}

type userController struct {
	userUseCase     usecase.UserUseCase
	customValidator helper.CustomValidator
	logs            *logger.Log
}

func NewUserController(userUseCase usecase.UserUseCase, customValidator helper.CustomValidator, logs *logger.Log) UserController {
	return &userController{userUseCase: userUseCase, customValidator: customValidator, logs: logs}
}

func (c *userController) GetUserProfile(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	userProfileResponse, err := c.userUseCase.GetUserProfile(ctx.Context(), auth.UserId)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.UserProfileResponse]{
		Success: true,
		Data:    userProfileResponse,
	})
}

func (c *userController) UpdateUserProfile(ctx *fiber.Ctx) error {
	request := new(model.RequestUpdateUserProfile)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	auth := middleware.GetUser(ctx)
	request.UserId = auth.UserId

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.userUseCase.UpdateUserProfile(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.UserProfileResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *userController) UpdateUserProfileImage(ctx *fiber.Ctx) error {

	userProfId := ctx.Params("userProfId", "")
	if userProfId == "" {
		return fiber.NewError(http.StatusBadRequest, "userProfId is required")
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "missing file: "+err.Error())
	}

	const maxFileSize = 5 * 1024 * 1024 // 5MB
	if file.Size > maxFileSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "File size exceeds the 2MB limit")
	}

	success, err := c.userUseCase.UpdateUserProfileImage(ctx.Context(), file, userProfId)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
		Data: map[string]interface{}{
			"success": success,
		},
	})
}

func (c *userController) UpdateUserCoverImage(ctx *fiber.Ctx) error {

	userProfId := ctx.Params("userProfId", "")
	if userProfId == "" {
		return fiber.NewError(http.StatusBadRequest, "userProfId is required")
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "missing file: "+err.Error())
	}

	const maxFileSize = 4 * 1024 * 1024 // 5MB
	if file.Size > maxFileSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "File size exceeds the 2MB limit")
	}

	success, err := c.userUseCase.UpdateUserCoverImage(ctx.Context(), file, userProfId)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
		Data: map[string]interface{}{
			"success": success,
		},
	})
}

func (c *userController) GetAllPublicUserChat(ctx *fiber.Ctx) error {
	log.Print("get all public accessed")
	request := &model.RequestGetAllPublicUser{
		Username: ctx.Query("username", ""),
		Page:     ctx.QueryInt("page", 1),
		Size:     ctx.QueryInt("size", 10),
	}

	response, pageMetadata, err := c.userUseCase.GetPublicUserChat(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, err, c.logs)
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*[]*model.GetAllPublicUserResponse]{
		Success:      true,
		Data:         response,
		PageMetadata: pageMetadata,
	})
}
