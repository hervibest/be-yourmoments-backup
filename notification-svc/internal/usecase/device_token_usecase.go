package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/repository"
	"github.com/oklog/ulid/v2"
)

type UserDeviceUseCase interface {
	CreateDevice(ctx context.Context, request *model.CreateDeviceRequest) error
}

type deviceTokenUseCase struct {
	db                   repository.BeginTx
	userDeviceRepository repository.UserDeviceRepository
	cacheAdapter         adapter.CacheAdapter
	logs                 logger.Log
}

func NewUserDeviceUseCase(db repository.BeginTx,
	userDeviceRepository repository.UserDeviceRepository,
	cacheAdapter adapter.CacheAdapter,
	logs logger.Log) UserDeviceUseCase {
	return &deviceTokenUseCase{
		db:                   db,
		userDeviceRepository: userDeviceRepository,
		cacheAdapter:         cacheAdapter,
		logs:                 logs,
	}
}

func (uc *deviceTokenUseCase) CreateDevice(ctx context.Context, request *model.CreateDeviceRequest) error {
	now := time.Now()
	userDevice := &entity.UserDevice{
		Id:        ulid.Make().String(),
		UserId:    request.UserID,
		Token:     request.DeviceToken,
		Platform:  request.Platform,
		CreatedAt: &now,
	}

	_, err := uc.userDeviceRepository.Create(ctx, uc.db, userDevice)
	if err != nil {
		return helper.WrapInternalServerError(uc.logs, "failed to create user device", err)
	}

	setKey := fmt.Sprintf("fcm_tokens:%s", userDevice.UserId)
	if err := uc.cacheAdapter.SAdd(ctx, setKey, request.DeviceToken); err != nil {
		return helper.WrapInternalServerError(uc.logs, "failed to SAdd redis set", err)
	}

	return nil
}
