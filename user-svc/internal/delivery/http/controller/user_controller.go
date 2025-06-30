package http

import (
	"log"
	"net/http"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	GetAllPublicUserChat(ctx *fiber.Ctx) error
	GetUserProfileV2(ctx *fiber.Ctx) error
	UpdateUserCoverImageV2(ctx *fiber.Ctx) error
	UpdateUserProfileImageV2(ctx *fiber.Ctx) error
	UpdateUserSimilarity(ctx *fiber.Ctx) error
}

type userController struct {
	userUseCase     usecase.UserUseCase
	customValidator helper.CustomValidator
	logs            logger.Log
}

func NewUserController(userUseCase usecase.UserUseCase, customValidator helper.CustomValidator, logs logger.Log) UserController {
	return &userController{userUseCase: userUseCase, customValidator: customValidator, logs: logs}
}

func (c *userController) GetUserProfileV2(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	userProfileResponse, err := c.userUseCase.GetUserProfileV2(ctx.Context(), auth.UserId)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Get user profile : ", err, c.logs)
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
		return helper.ErrUseCaseResponseJSON(ctx, "Update user profile : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.UserProfileResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *userController) UpdateUserProfileImageV2(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "missing file: "+err.Error())
	}

	const maxFileSize = 5 * 1024 * 1024 // 5MB
	if file.Size > maxFileSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "File size exceeds the 2MB limit")
	}

	c.logs.Log("Update user profile iuamge v2 accessed")
	auth := middleware.GetUser(ctx)
	success, err := c.userUseCase.UpdateUserCoverImageV2(ctx.Context(), file, auth.UserProfileID)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Update user profile image : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
		Data: map[string]interface{}{
			"success": success,
		},
	})
}

func (c *userController) UpdateUserCoverImageV2(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "missing file: "+err.Error())
	}

	const maxFileSize = 5 * 1024 * 1024 // 5MB
	if file.Size > maxFileSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "File size exceeds the 5MB limit")
	}

	auth := middleware.GetUser(ctx)
	success, err := c.userUseCase.UpdateUserCoverImageV2(ctx.Context(), file, auth.UserProfileID)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Update user cover image : ", err, c.logs)
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
		return helper.ErrUseCaseResponseJSON(ctx, "Get all public user chat : ", err, c.logs)
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*[]*model.GetAllPublicUserResponse]{
		Success:      true,
		Data:         response,
		PageMetadata: pageMetadata,
	})
}

func (c *userController) UpdateUserSimilarity(ctx *fiber.Ctx) error {
	request := new(model.RequestUpdateSimilarity)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	auth := middleware.GetUser(ctx)
	request.UserID = auth.UserId

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.userUseCase.UpdateUserSimilarity(ctx.Context(), request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Update user profile : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}
