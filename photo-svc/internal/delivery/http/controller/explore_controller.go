package http

import (
	"be-yourmoments/photo-svc/internal/delivery/http/middleware"
	"be-yourmoments/photo-svc/internal/helper"
	"be-yourmoments/photo-svc/internal/helper/logger"
	"be-yourmoments/photo-svc/internal/model"
	"be-yourmoments/photo-svc/internal/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type ExploreController interface {
	GetAllExploreSimilar(ctx *fiber.Ctx) error
	GetAllUserCart(ctx *fiber.Ctx) error
	GetAllUserFavorite(ctx *fiber.Ctx) error
	GetAllUserWishlist(ctx *fiber.Ctx) error
	UserAddCart(ctx *fiber.Ctx) error
	UserAddFavorite(ctx *fiber.Ctx) error
	UserAddWishlist(ctx *fiber.Ctx) error
	UserDeleteCart(ctx *fiber.Ctx) error
	UserDeleteFavorite(ctx *fiber.Ctx) error
	UserDeleteWishlist(ctx *fiber.Ctx) error
}

type exploreController struct {
	exploreUseCase  usecase.ExploreUseCase
	customValidator helper.CustomValidator
	tracer          trace.Tracer
	logs            *logger.Log
}

func NewExploreController(tracer trace.Tracer, customValidator helper.CustomValidator, exploreUseCase usecase.ExploreUseCase,
	logs *logger.Log) ExploreController {
	return &exploreController{exploreUseCase: exploreUseCase, customValidator: customValidator, tracer: tracer, logs: logs}
}

func (c *exploreController) GetAllExploreSimilar(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	context, span := c.tracer.Start(ctx.Context(), "getUser", oteltrace.WithAttributes(attribute.String("id", auth.UserId)))
	defer span.End()
	request := &model.GetAllExploreSimilarRequest{
		UserId: auth.UserId,
		Page:   ctx.QueryInt("page", 1),
		Size:   ctx.QueryInt("size", 10),
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, pageMetadata, err := c.exploreUseCase.GetUserExploreSimilar(context, request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Get all explore similar error : ", err, c.logs)
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*[]*model.ExploreUserSimilarResponse]{
		Success:      true,
		Data:         response,
		PageMetadata: pageMetadata,
	})
}

func (c *exploreController) GetAllUserWishlist(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	context, span := c.tracer.Start(ctx.Context(), "getUserWishlist", oteltrace.WithAttributes(attribute.String("id", auth.UserId)))
	defer span.End()
	request := &model.GetAllWishlistRequest{
		UserId: auth.UserId,
		Page:   ctx.QueryInt("page", 1),
		Size:   ctx.QueryInt("size", 10),
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, pageMetadata, err := c.exploreUseCase.GetUserWishlist(context, request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Get all user error : ", err, c.logs)
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*[]*model.ExploreUserSimilarResponse]{
		Success:      true,
		Data:         response,
		PageMetadata: pageMetadata,
	})
}

func (c *exploreController) UserAddWishlist(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	context, span := c.tracer.Start(ctx.Context(), "userAddWishlist", oteltrace.WithAttributes(attribute.String("id", auth.UserId)))
	defer span.End()

	request := &model.UserAddWishlistRequest{
		UserId: auth.UserId,
	}

	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if _, err := ulid.Parse(request.PhotoId); err != nil {
		return fiber.NewError(http.StatusBadRequest, "The provided photo ID is not valid")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.exploreUseCase.UserAddWishlist(context, request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "User add wishlist error : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *exploreController) UserDeleteWishlist(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	context, span := c.tracer.Start(ctx.Context(), "userAddWishlist", oteltrace.WithAttributes(attribute.String("id", auth.UserId)))
	defer span.End()

	request := &model.UserDeleteWishlistReqeust{
		UserId: auth.UserId,
	}

	if _, err := ulid.Parse(request.PhotoId); err != nil {
		return fiber.NewError(http.StatusBadRequest, "The provided photo ID is not valid")
	}

	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.exploreUseCase.UserDeleteWishlist(context, request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "User delete wishlist error : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *exploreController) GetAllUserFavorite(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	context, span := c.tracer.Start(ctx.Context(), "getUserFavorite", oteltrace.WithAttributes(attribute.String("id", auth.UserId)))
	defer span.End()
	request := &model.GetAllFavoriteRequest{
		UserId: auth.UserId,
		Page:   ctx.QueryInt("page", 1),
		Size:   ctx.QueryInt("size", 10),
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, pageMetadata, err := c.exploreUseCase.GetUserFavorite(context, request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "User get all user favorite error : ", err, c.logs)
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*[]*model.ExploreUserSimilarResponse]{
		Success:      true,
		Data:         response,
		PageMetadata: pageMetadata,
	})
}

func (c *exploreController) UserAddFavorite(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	context, span := c.tracer.Start(ctx.Context(), "userAddFavorite", oteltrace.WithAttributes(attribute.String("id", auth.UserId)))
	defer span.End()

	request := &model.UserAddFavoriteRequest{
		UserId: auth.UserId,
	}

	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if _, err := ulid.Parse(request.PhotoId); err != nil {
		return fiber.NewError(http.StatusBadRequest, "The provided photo ID is not valid")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.exploreUseCase.UserAddFavorite(context, request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "User add favorite error : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *exploreController) UserDeleteFavorite(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	context, span := c.tracer.Start(ctx.Context(), "userAddFavorite", oteltrace.WithAttributes(attribute.String("id", auth.UserId)))
	defer span.End()

	request := &model.UserDeleteFavoriteReqeust{
		UserId: auth.UserId,
	}

	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if _, err := ulid.Parse(request.PhotoId); err != nil {
		return fiber.NewError(http.StatusBadRequest, "The provided photo ID is not valid")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.exploreUseCase.UserDeleteFavorite(context, request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "User delete favorite error : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *exploreController) GetAllUserCart(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	context, span := c.tracer.Start(ctx.Context(), "getUserCart", oteltrace.WithAttributes(attribute.String("id", auth.UserId)))
	defer span.End()
	request := &model.GetAllCartRequest{
		UserId: auth.UserId,
		Page:   ctx.QueryInt("page", 1),
		Size:   ctx.QueryInt("size", 10),
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	response, pageMetadata, err := c.exploreUseCase.GetUserCart(context, request)
	if err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "Get all user cart error : ", err, c.logs)
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*[]*model.ExploreUserSimilarResponse]{
		Success:      true,
		Data:         response,
		PageMetadata: pageMetadata,
	})
}

func (c *exploreController) UserAddCart(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	context, span := c.tracer.Start(ctx.Context(), "userAddCart", oteltrace.WithAttributes(attribute.String("id", auth.UserId)))
	defer span.End()

	request := &model.UserAddCartRequest{
		UserId: auth.UserId,
	}

	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if _, err := ulid.Parse(request.PhotoId); err != nil {
		return fiber.NewError(http.StatusBadRequest, "The provided photo ID is not valid")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.exploreUseCase.UserAddCart(context, request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "User add cart error : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *exploreController) UserDeleteCart(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	context, span := c.tracer.Start(ctx.Context(), "userDeleteCart", oteltrace.WithAttributes(attribute.String("id", auth.UserId)))
	defer span.End()

	request := &model.UserDeleteCartReqeust{
		UserId: auth.UserId,
	}

	if err := helper.StrictBodyParser(ctx, request); err != nil {
		return helper.ErrBodyParserResponseJSON(ctx, err)
	}

	if _, err := ulid.Parse(request.PhotoId); err != nil {
		return fiber.NewError(http.StatusBadRequest, "The provided photo ID is not valid")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return helper.ErrValidationResponseJSON(ctx, validatonErrs)
	}

	if err := c.exploreUseCase.UserDeleteCart(context, request); err != nil {
		return helper.ErrUseCaseResponseJSON(ctx, "User delete cart : ", err, c.logs)
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}
