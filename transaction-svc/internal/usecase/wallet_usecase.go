package usecase

import (
	"be-yourmoments/transaction-svc/internal/entity"
	"be-yourmoments/transaction-svc/internal/helper"
	"be-yourmoments/transaction-svc/internal/helper/logger"
	"be-yourmoments/transaction-svc/internal/model"
	"be-yourmoments/transaction-svc/internal/model/converter"
	"be-yourmoments/transaction-svc/internal/repository"
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type WalletUsecase interface {
	CreateWallet(ctx context.Context, request *model.RequestCreateWallet) (*model.WalletResponse, error)
}
type walletUsecase struct {
	walletRepository repository.WalletRepository
	db               *sqlx.DB
	logs             *logger.Log
}

func NewWalletUsecase(walletRepository repository.WalletRepository, db *sqlx.DB, logs *logger.Log) WalletUsecase {
	log.Printf("wallet usecase initialized")
	return &walletUsecase{walletRepository: walletRepository, db: db, logs: logs}
}

func (u *walletUsecase) CreateWallet(ctx context.Context, request *model.RequestCreateWallet) (*model.WalletResponse, error) {
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
