package http

import (
	"be-yourmoments/photo-svc/internal/delivery/http/middleware"
	"be-yourmoments/photo-svc/internal/helper"
	"be-yourmoments/photo-svc/internal/helper/logger"
	"be-yourmoments/photo-svc/internal/model"
	"be-yourmoments/photo-svc/internal/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type CheckoutController interface {
	PreviewCheckout(ctx *fiber.Ctx) error
}
type checkoutController struct {
	checkoutUseCase usecase.CheckoutUseCase
	customValidator helper.CustomValidator
	logs            *logger.Log
}

func NewCheckoutController(checkoutUseCase usecase.CheckoutUseCase, customValidator helper.CustomValidator, logs *logger.Log) CheckoutController {
	return &checkoutController{checkoutUseCase: checkoutUseCase, customValidator: customValidator, logs: logs}
}

func (c *checkoutController) PreviewCheckout(ctx *fiber.Ctx) error {
	request := new(model.PreviewCheckoutRequest)
	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	user := middleware.GetUser(ctx)
	request.UserId = user.UserId

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, err := c.checkoutUseCase.PreviewCheckout(ctx.Context(), request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Preview checkout : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.PreviewCheckoutResponse]{
		Success: true,
		Data:    response,
	})
}
