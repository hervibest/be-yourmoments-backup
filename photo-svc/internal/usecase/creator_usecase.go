package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
	errorcode "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/converter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/repository"

	"github.com/oklog/ulid/v2"
)

type CreatorUseCase interface {
	CreateCreator(ctx context.Context, req *model.CreateCreatorRequest) (*model.CreatorResponse, error)
	GetCreator(ctx context.Context, req *model.GetCreatorRequest) (*model.CreatorResponse, error)
}

type creatorUseCase struct {
	db                 repository.BeginTx
	creatorRepository  repository.CreatorRepository
	transactionAdapter adapter.TransactionAdapter
	logs               *logger.Log
}

func NewCreatorUseCase(db repository.BeginTx, creatorRepository repository.CreatorRepository, transactionAdapter adapter.TransactionAdapter, logs *logger.Log) CreatorUseCase {
	return &creatorUseCase{
		db:                 db,
		creatorRepository:  creatorRepository,
		transactionAdapter: transactionAdapter,
		logs:               logs}
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

	// _, err = u.transactionAdapter.CreateWallet(ctx, creator.Id)
	// if err != nil {
	// 	return nil, err
	// }

	return converter.CreatorToResponse(creator), nil
}

func (u *creatorUseCase) GetCreator(ctx context.Context, request *model.GetCreatorRequest) (*model.CreatorResponse, error) {
	creator, err := u.creatorRepository.FindByUserId(ctx, request.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid user id")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to find creator by user id", err)
	}

	return converter.CreatorToResponse(creator), nil
}
