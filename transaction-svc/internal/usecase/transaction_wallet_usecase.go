package usecase

import (
	"be-yourmoments/transaction-svc/internal/helper"
	"be-yourmoments/transaction-svc/internal/helper/logger"
	"be-yourmoments/transaction-svc/internal/model"
	"be-yourmoments/transaction-svc/internal/model/converter"
	"be-yourmoments/transaction-svc/internal/repository"
	"context"

	"github.com/jmoiron/sqlx"
)

type TransactionWalletUseCase interface {
	GetAll(ctx context.Context, request *model.GetAllTransactionWallet) (*[]*model.TransactionWalletResponse, *model.PageMetadata, error)
}

type transactionWalletUseCase struct {
	db                    *sqlx.DB
	transactionWalletRepo repository.TransactionWalletRepository
	logs                  *logger.Log
}

func NewTransactionWalletUseCase(db *sqlx.DB, transactionWalletRepo repository.TransactionWalletRepository, logs *logger.Log) TransactionWalletUseCase {
	return &transactionWalletUseCase{db: db, transactionWalletRepo: transactionWalletRepo, logs: logs}
}

func (u *transactionWalletUseCase) GetAll(ctx context.Context, request *model.GetAllTransactionWallet) (*[]*model.TransactionWalletResponse, *model.PageMetadata, error) {
	wallets, pagination, err := u.transactionWalletRepo.FindAllByWalletId(ctx, u.db, request.Page, request.Size,
		request.WalletId, request.Max, request.Min, request.Order)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to distribute transaction to wallets", err)
	}

	return converter.WalletsToResponses(&wallets), pagination, nil
}
