package middleware

import (
	"errors"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TODO SHOULD TOKEN VALIDATED ?
func NewUserAuth(userAdapter adapter.UserAdapter, tracer trace.Tracer, logs *logger.Log) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		context, span := tracer.Start(ctx.Context(), "authenticateUser", oteltrace.WithAttributes())
		defer span.End()
		token := strings.TrimPrefix(ctx.Get("Authorization", ""), "Bearer ")
		if token == "" || token == "NOT_FOUND" {
			return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized access")
		}

		authResponse, err := userAdapter.AuthenticateUser(context, token)
		if err != nil {
			return helper.ErrUseCaseResponseJSON(ctx, "Authenticate user error : ", err, logs)
		}

		auth := &model.AuthResponse{
			UserId:      authResponse.GetUser().GetUserId(),
			Username:    authResponse.GetUser().GetUsername(),
			Email:       authResponse.GetUser().GetEmail(),
			PhoneNumber: authResponse.GetUser().GetPhoneNumber(),
			Similarity:  authResponse.GetUser().GetSimilarity(),
			// CreatorId:   authResponse.GetUser().GetCreatorId(),
			// WalletId:    authResponse.GetUser().GetWalletId(),
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

type AuthMiddleware interface {
	NewUserAuthV2() fiber.Handler
}
type authMiddleware struct {
	userAdapter  adapter.UserAdapter
	logs         *logger.Log
	tracer       trace.Tracer
	jwtAdapter   adapter.JWTAdapter
	cacheAdapter adapter.CacheAdapter
}

func NewAuthMiddleware(userAdapter adapter.UserAdapter, logs *logger.Log, tracer trace.Tracer,
	jwtAdapter adapter.JWTAdapter, cacheAdapter adapter.CacheAdapter) AuthMiddleware {
	return &authMiddleware{
		userAdapter:  userAdapter,
		logs:         logs,
		tracer:       tracer,
		jwtAdapter:   jwtAdapter,
		cacheAdapter: cacheAdapter,
	}
}

func (m *authMiddleware) NewUserAuthV2() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		m.logs.Log("Accessed the user auth new middleware")
		context, span := m.tracer.Start(ctx.Context(), "authenticateUser", oteltrace.WithAttributes())
		defer span.End()
		token := strings.TrimPrefix(ctx.Get("Authorization", ""), "Bearer ")
		if token == "" || token == "NOT_FOUND" {
			return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized access")
		}

		accessTokenDetail, err := m.jwtAdapter.VerifyAccessToken(token)
		if err != nil {
			return fiber.NewError(fiber.ErrUnauthorized.Code, "Invalid access token")
		}

		userId, _ := m.cacheAdapter.Get(context, token)
		if userId != "" {
			return fiber.NewError(fiber.ErrUnauthorized.Code, "User already signed out")
		}

		cachedUserStr, err := m.cacheAdapter.Get(context, accessTokenDetail.UserId)
		if err != nil && !errors.Is(err, redis.Nil) {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong, please try again later")
		}

		authEntity := new(entity.Auth)
		//If redis stale, get from db
		if errors.Is(err, redis.Nil) {
			authResponse, err := m.userAdapter.AuthenticateUserV2(context, token, accessTokenDetail.UserId, accessTokenDetail.ExpiresAt)
			if err != nil {
				return fiber.NewError(fiber.ErrBadRequest.Code, "Invalid access token")
			}
			authEntity.Email = authResponse.GetUser().GetEmail()
			authEntity.Id = authResponse.GetUser().GetUserId()
			authEntity.PhoneNumber = authResponse.GetUser().GetPhoneNumber()
			authEntity.Similarity = uint(authResponse.GetUser().GetSimilarity())

		} else {
			if err := sonic.ConfigFastest.Unmarshal([]byte(cachedUserStr), &authEntity); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong, please try again later")
			}
		}

		auth := &model.AuthResponse{
			UserId:      authEntity.Id,
			Username:    authEntity.Username,
			Email:       authEntity.Email,
			PhoneNumber: authEntity.PhoneNumber,
			Similarity:  uint32(authEntity.Similarity),
			// CreatorId:   authResponse.GetUser().GetCreatorId(),
			// WalletId:    authResponse.GetUser().GetWalletId(),
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.AuthResponse {
	return ctx.Locals("auth").(*model.AuthResponse)
}
