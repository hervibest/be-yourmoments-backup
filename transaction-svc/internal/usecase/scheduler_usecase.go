package usecase

import (
	"context"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/enum"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/repository"
	"github.com/jmoiron/sqlx"
)

type SchedulerUseCase interface{}
type schedulerUseCase struct {
	db                    *sqlx.DB
	transactionRepository repository.TransactionRepository
	paymentAdapter        adapter.PaymentAdapter
	logs                  *logger.Log
}

func NewSchedulerUseCase(db *sqlx.DB, transactionRepository repository.TransactionRepository, paymentAdapter adapter.PaymentAdapter, logs *logger.Log) SchedulerUseCase {
	return &schedulerUseCase{db: db, transactionRepository: transactionRepository, paymentAdapter: paymentAdapter, logs: logs}
}

func (u *schedulerUseCase) CheckTransactionStatus(ctx context.Context) error {
	transactions, err := u.transactionRepository.FindManyByStatus(ctx, u.db, enum.TransactionStatusPending)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to find many transaction by status", err)
	}

	for _, transaction := range *transactions {
		_, err := u.paymentAdapter.CheckTransactionStatus(ctx, transaction.Id)
		if err != nil {
			return err
		}

	}

	return nil
}
