package http

import (
	"net/http"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type ChatController interface {
	GetCustomToken(ctx *fiber.Ctx) error
	GetOrCreateRoom(ctx *fiber.Ctx) error
	SendMessage(ctx *fiber.Ctx) error
}

type chatController struct {
	chatUseCase     usecase.ChatUseCase
	customValidator helper.CustomValidator
	logs            *logger.Log
}

func NewChatController(chatUseCase usecase.ChatUseCase, customValidator helper.CustomValidator, logs *logger.Log) ChatController {
	return &chatController{chatUseCase: chatUseCase, customValidator: customValidator, logs: logs}
}

func (c *chatController) GetCustomToken(ctx *fiber.Ctx) error {
	request := new(model.RequestCustomToken)
	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.chatUseCase.GetCustomToken(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Get custom token error : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.CustomTokenResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *chatController) GetOrCreateRoom(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.RequestGetOrCreateRoom{
		SenderId: auth.UserId,
	}

	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.chatUseCase.GetOrCreateRoom(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Get custom or create room error : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.GetOrCreateRoomResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *chatController) SendMessage(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.RequestSendMessage{
		SenderId: auth.UserId,
	}

	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.chatUseCase.SendMessage(ctx.Context(), request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Send message : ", err, c.logs)
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.GetOrCreateRoomResponse]{
		Success: true,
	})
}
