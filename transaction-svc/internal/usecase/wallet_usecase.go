package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	errorcode "github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model/converter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/repository"
	"github.com/redis/go-redis/v9"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type WalletUseCase interface {
	CreateWallet(ctx context.Context, request *model.CreateWalletRequest) (*model.WalletResponse, error)
	GetWallet(ctx context.Context, request *model.GetWalletRequest) (*model.WalletResponse, error)
	GetWalletId(ctx context.Context, request *model.GetWalletIdRequest) (string, error)
}
type walletUseCase struct {
	walletRepository repository.WalletRepository
	cacheAdapter     adapter.CacheAdapter
	db               *sqlx.DB
	logs             *logger.Log
}

func NewWalletUseCase(walletRepository repository.WalletRepository, cacheAdapter adapter.CacheAdapter,
	db *sqlx.DB, logs *logger.Log) WalletUseCase {
	log.Printf("wallet usecase initialized")
	return &walletUseCase{walletRepository: walletRepository, cacheAdapter: cacheAdapter, db: db, logs: logs}
}

func (u *walletUseCase) CreateWallet(ctx context.Context, request *model.CreateWalletRequest) (*model.WalletResponse, error) {
	now := time.Now()
	wallet := &entity.Wallet{
		Id:        ulid.Make().String(),
		CreatorId: request.CreatorId,
		Balance:   0,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return nil, err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	wallet, err = u.walletRepository.Create(ctx, tx, wallet)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to create wallet in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return nil, err
	}

	return converter.WalletToResponse(wallet), nil
}

func (u *walletUseCase) GetWallet(ctx context.Context, request *model.GetWalletRequest) (*model.WalletResponse, error) {
	wallet, err := u.walletRepository.FindByCreatorId(ctx, u.db, request.CreatorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.NewUseCaseError(errorcode.ErrResourceNotFound, "Wallet not found make sure to give a valid creator id")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to find wallet by creator id", err)
	}

	return converter.WalletToResponse(wallet), nil
}

func (u *walletUseCase) GetWalletId(ctx context.Context, request *model.GetWalletIdRequest) (string, error) {
	walletId, err := u.cacheAdapter.Get(ctx, request.CreatorId)
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", helper.WrapInternalServerError(u.logs, "failed to get cached user", err)
	}

	if errors.Is(err, redis.Nil) {
		walletId, err = u.walletRepository.FindIdByCreatorId(ctx, u.db, request.CreatorId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", helper.NewUseCaseError(errorcode.ErrResourceNotFound, "Wallet not found make sure to give a valid creator id")
			}
			return "", helper.WrapInternalServerError(u.logs, "failed to find wallet by creator id", err)
		}

		if err := u.cacheAdapter.Set(ctx, request.CreatorId, walletId, 240*time.Minute); err != nil {
			return "", helper.WrapInternalServerError(u.logs, "failed to save wallet id to cache", err)
		}
	}

	return walletId, nil
}
