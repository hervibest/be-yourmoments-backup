package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/bytedance/sonic"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
	errorcode "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum/error"
	producer "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/gateway/messaging"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/converter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/event"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/repository"
	"github.com/jmoiron/sqlx"

	"github.com/redis/go-redis/v9"

	"github.com/oklog/ulid/v2"
)

type CreatorUseCase interface {
	CreateCreator(ctx context.Context, req *model.CreateCreatorRequest) (*model.CreatorResponse, error)
	GetCreator(ctx context.Context, req *model.GetCreatorRequest) (*model.CreatorResponse, error)
	UpdateCreatorTotalReview(ctx context.Context, req *model.UpdateCreatorTotalRatingRequest) (*model.CreatorResponse, error)
	GetCreatorId(ctx context.Context, request *model.GetCreatorIdRequest) (string, error)
}

type creatorUseCase struct {
	db                *sqlx.DB
	creatorRepository repository.CreatorRepository
	cacheAdapter      adapter.CacheAdapter
	creatorProducer   producer.CreatorProducer
	logs              *logger.Log
}

func NewCreatorUseCase(db *sqlx.DB, creatorRepository repository.CreatorRepository, cacheAdapter adapter.CacheAdapter,
	creatorProducer producer.CreatorProducer, logs *logger.Log) CreatorUseCase {
	return &creatorUseCase{
		db:                db,
		creatorRepository: creatorRepository,
		cacheAdapter:      cacheAdapter,
		creatorProducer:   creatorProducer,
		logs:              logs,
	}
}

func (u *creatorUseCase) CreateCreator(ctx context.Context, req *model.CreateCreatorRequest) (*model.CreatorResponse, error) {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return nil, err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	now := time.Now()

	creator := &entity.Creator{
		Id:        ulid.Make().String(),
		UserId:    req.UserId,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	creator, err = u.creatorRepository.Create(ctx, tx, creator)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to insert new creator to database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return nil, err
	}

	event := &event.CreatorEvent{
		Id:        creator.Id,
		UserId:    creator.UserId,
		CreatedAt: creator.CreatedAt,
		UpdatedAt: creator.UpdatedAt,
	}

	if err := u.creatorProducer.ProduceCreatorCreated(ctx, event); err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to procuer creator created", err)
	}

	return converter.CreatorToResponse(creator), nil
}

func (u *creatorUseCase) GetCreator(ctx context.Context, request *model.GetCreatorRequest) (*model.CreatorResponse, error) {
	creator := new(entity.Creator)
	creatorJson, err := u.cacheAdapter.Get(ctx, "creator:"+request.UserId)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, helper.WrapInternalServerError(u.logs, "failed to get cached user", err)
	}

	if errors.Is(err, redis.Nil) {
		creator, err = u.creatorRepository.FindByUserId(ctx, request.UserId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid user id")
			}
			return nil, helper.WrapInternalServerError(u.logs, "failed to find creator by user id", err)
		}

		creatorByte, err := sonic.ConfigFastest.Marshal(creator)
		if err != nil {
			return nil, helper.WrapInternalServerError(u.logs, "failed to marshal creator", err)
		}

		if err := u.cacheAdapter.Set(ctx, "creator:"+request.UserId, creatorByte, 240*time.Minute); err != nil {
			return nil, helper.WrapInternalServerError(u.logs, "failed to save creator to cache", err)
		}
	} else {
		if err := sonic.ConfigFastest.Unmarshal([]byte(creatorJson), creator); err != nil {
			return nil, helper.WrapInternalServerError(u.logs, "failed to unmarshal creator", err)
		}
	}

	return converter.CreatorToResponse(creator), nil
}

func (u *creatorUseCase) GetCreatorId(ctx context.Context, request *model.GetCreatorIdRequest) (string, error) {
	creatorId, err := u.cacheAdapter.Get(ctx, request.UserId)
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", helper.WrapInternalServerError(u.logs, "failed to get cached user", err)
	}

	if errors.Is(err, redis.Nil) {
		creatorId, err = u.creatorRepository.FindIdByUserId(ctx, u.db, request.UserId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", helper.NewUseCaseError(errorcode.ErrResourceNotFound, "Creator not found make sure to give a valid creator id")
			}
			return "", helper.WrapInternalServerError(u.logs, "failed to find wallet by creator id", err)
		}

		if err := u.cacheAdapter.Set(ctx, request.UserId, creatorId, 240*time.Minute); err != nil {
			return "", helper.WrapInternalServerError(u.logs, "failed to save creator to cache", err)
		}
	}

	return creatorId, nil
}

func (u *creatorUseCase) UpdateCreatorTotalReview(ctx context.Context, req *model.UpdateCreatorTotalRatingRequest) (*model.CreatorResponse, error) {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return nil, err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	now := time.Now()

	creator := &entity.Creator{
		Id:          req.Id,
		Rating:      req.Rating,
		RatingCount: req.RatingCount,
		UpdatedAt:   &now,
	}

	creator, err = u.creatorRepository.UpdateCreatorRating(ctx, tx, creator)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to update creator total review to database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return nil, err
	}

	creatorByte, err := sonic.ConfigFastest.Marshal(creator)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to marshal creator", err)
	}

	if err := u.cacheAdapter.Set(ctx, "creator:"+creator.UserId, creatorByte, 240*time.Minute); err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to save creator to cache", err)
	}

	return converter.CreatorToResponse(creator), nil
}
