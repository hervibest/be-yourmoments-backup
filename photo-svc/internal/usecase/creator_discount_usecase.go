package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
	errorcode "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/converter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type CreatorDiscountUseCase interface {
	CreateDiscount(ctx context.Context, request *model.CreateCreatorDiscountRequest) (*model.CreatorDiscountResponse, error)
	ActivateDiscount(ctx context.Context, request *model.ActivateCreatorDiscountRequest) error
	DeactivateDiscount(ctx context.Context, request *model.DeactivateCreatorDiscountRequest) error
	GetDiscount(ctx context.Context, request *model.GetCreatorDiscountRequest) (*model.CreatorDiscountResponse, error)
	GetAllDiscount(ctx context.Context, creatorId string) (*[]*model.CreatorDiscountResponse, error)
}
type creatorDiscountUseCase struct {
	db                        *sqlx.DB
	creatorDiscountRepository repository.CreatorDiscountRepository
	logs                      *logger.Log
}

func NewCreatorDiscountUseCase(db *sqlx.DB, creatorDiscountRepository repository.CreatorDiscountRepository, logs *logger.Log) CreatorDiscountUseCase {
	return &creatorDiscountUseCase{
		db:                        db,
		creatorDiscountRepository: creatorDiscountRepository,
		logs:                      logs,
	}
}

// TODO validate price when there is a photo price
func (u *creatorDiscountUseCase) CreateDiscount(ctx context.Context, request *model.CreateCreatorDiscountRequest) (*model.CreatorDiscountResponse, error) {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return nil, err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if request.DiscountType != enum.DiscountTypeFlat && request.DiscountType != enum.DiscountTypePercent {
		return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid discount type")
	}

	now := time.Now()

	creatorDiscount := &entity.CreatorDiscount{
		Id:           ulid.Make().String(),
		CreatorId:    request.CreatorId,
		Name:         request.Name,
		MinQuantity:  request.MinQuantity,
		DiscountType: request.DiscountType,
		Value:        request.Value,
		IsActive:     request.IsActive,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}

	u.logs.CustomLog("create discount creator ", creatorDiscount)
	creatorDiscount, err = u.creatorDiscountRepository.Create(ctx, tx, creatorDiscount)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to create creator discount to database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return nil, err
	}

	return converter.CreatorDiscountToResponse(creatorDiscount), nil
}

func (u *creatorDiscountUseCase) ActivateDiscount(ctx context.Context, request *model.ActivateCreatorDiscountRequest) error {
	_, err := u.creatorDiscountRepository.FindByIdAndCreatorId(ctx, u.db, request.Id, request.CreatorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewUseCaseError(errorcode.ErrResourceNotFound, "Invalid creator discount id")
		}
		return helper.WrapInternalServerError(u.logs, "failed to find creator discount by discount id", err)
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if err = u.creatorDiscountRepository.Activate(ctx, tx, request.Id); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to activate creator discount in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}

func (u *creatorDiscountUseCase) DeactivateDiscount(ctx context.Context, request *model.DeactivateCreatorDiscountRequest) error {
	_, err := u.creatorDiscountRepository.FindByIdAndCreatorId(ctx, u.db, request.Id, request.CreatorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid creator discount id")
		}
		return helper.WrapInternalServerError(u.logs, "failed to find creator discount by discount id", err)
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if err = u.creatorDiscountRepository.Deactivate(ctx, tx, request.Id); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to deactivate creator discount in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}

func (u *creatorDiscountUseCase) GetDiscount(ctx context.Context, request *model.GetCreatorDiscountRequest) (*model.CreatorDiscountResponse, error) {
	discount, err := u.creatorDiscountRepository.FindByIdAndCreatorId(ctx, u.db, request.Id, request.CreatorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.NewUseCaseError(errorcode.ErrResourceNotFound, "Creator discount not found")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to find creator discount by discount id", err)
	}

	return converter.CreatorDiscountToResponse(discount), nil
}

func (u *creatorDiscountUseCase) GetAllDiscount(ctx context.Context, creatorId string) (*[]*model.CreatorDiscountResponse, error) {
	discounts, err := u.creatorDiscountRepository.FindAll(ctx, u.db, creatorId)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to find  all creator discount", err)
	}
	return converter.CreatorDiscountsToResponses(*discounts), nil
}
