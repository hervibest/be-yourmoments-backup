package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/enum"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/repository"
	"github.com/jmoiron/sqlx"
)

type CancelationUseCase interface {
	ExpirePendingTransaction(ctx context.Context, transactionId string) error
	CancelPendingTransaction(ctx context.Context, transactionId string) error
}

type cancelationUseCase struct {
	db              *sqlx.DB
	transactionRepo repository.TransactionRepository
	logs            *logger.Log
}

func NewCancelationUseCase(db *sqlx.DB, transactionRepo repository.TransactionRepository, logs *logger.Log) CancelationUseCase {
	return &cancelationUseCase{db: db, transactionRepo: transactionRepo, logs: logs}
}

func (u *cancelationUseCase) ExpirePendingTransaction(ctx context.Context, transactionId string) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}
	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if err := u.updateTransactionStatusIfPending(ctx, tx, transactionId, enum.TrxInternalStatusExpired); err != nil {
		return helper.WrapInternalServerError(u.logs, "expire failed", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}
	u.logs.Log(fmt.Sprintf("success expire pending transaction with id : %s", transactionId))
	return nil
}

func (u *cancelationUseCase) CancelPendingTransaction(ctx context.Context, transactionId string) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}
	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if err := u.updateTransactionStatusIfPending(ctx, tx, transactionId, enum.TrxInternalStatusCancelledByUser); err != nil {
		return helper.WrapInternalServerError(u.logs, "cancel failed", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}
	u.logs.Log(fmt.Sprintf("success cancel pending by user transaction with id : %s", transactionId))
	return nil
}

func (u *cancelationUseCase) updateTransactionStatusIfPending(ctx context.Context, tx *sqlx.Tx, transactionId string, status enum.TrxInternalStatus) error {
	transaction, err := u.transactionRepo.FindById(ctx, tx, transactionId, true)
	if err != nil {
		return err
	}

	if transaction.InternalStatus != enum.TrxInternalStatusPending && transaction.InternalStatus != enum.TrxInternalStatusTokenReady {
		u.logs.Log(fmt.Sprintf("transaction : %s is not in pending or token ready status, discontinued update transaction process status if pending", transactionId))
		return nil
	}

	now := time.Now()
	if status == enum.TrxInternalStatusExpired {
		transaction.Status = enum.TransactionStatusExpired
		transaction.InternalStatus = enum.TrxInternalStatusExpired
	} else {
		transaction.Status = enum.TransactionStatusCancelled
		transaction.InternalStatus = enum.TrxInternalStatusCancelledByUser
	}
	transaction.SnapToken = sql.NullString{Valid: true, String: ""}
	transaction.UpdatedAt = &now

	return u.transactionRepo.UpdateStatus(ctx, tx, transaction)
}
