package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/enum"
	errorcode "github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model/converter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type WithdrawalUseCase interface {
	Create(ctx context.Context, request *model.CreateWithdrawalRequest) (*model.WithdrawalResponse, error)
	FindAll(ctx context.Context) (*[]*model.WithdrawalResponse, error)
	FindById(ctx context.Context, request *model.FindWithdrawalById) (*model.WithdrawalResponse, error)
	Update(ctx context.Context, request *model.UpdateWithdrawalStatusRequest) (*model.WithdrawalResponse, error)
}

type withdrawalUseCase struct {
	db                   *sqlx.DB
	withdrawalRepository repository.WithdrawalRepository
	walletRepository     repository.WalletRepository
	logs                 *logger.Log
}

func NewWithdrawalUseCase(db *sqlx.DB, withdrawalRepository repository.WithdrawalRepository, walletRepository repository.WalletRepository,
	logs *logger.Log) WithdrawalUseCase {
	return &withdrawalUseCase{
		db:                   db,
		withdrawalRepository: withdrawalRepository,
		walletRepository:     walletRepository,
		logs:                 logs,
	}
}

func (u *withdrawalUseCase) Create(ctx context.Context, request *model.CreateWithdrawalRequest) (*model.WithdrawalResponse, error) {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return nil, err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	wallet, err := u.walletRepository.FindById(ctx, tx, request.WalletId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.NewUseCaseError(errorcode.ErrResourceNotFound, "Invalid creator discount id")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to find creator discount by discount id", err)
	}

	if wallet.Balance-int32(request.Amount) < 0 {
		return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Negative amount")
	}

	// TODO reduce balance
	if err := u.walletRepository.ReduceBalance(ctx, tx, request.WalletId, int64(request.Amount)); err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to find creator discount by discount id", err)
	}

	now := time.Now()
	withdrawal := &entity.Withdrawal{
		Id:           ulid.Make().String(),
		WalletId:     wallet.Id,
		BankWalletId: request.BankWalletId,
		Amount:       request.Amount,
		Status:       enum.WithdrawalStatusPending,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}

	withdrawal, err = u.withdrawalRepository.Create(ctx, tx, withdrawal)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to find creator discount by discount id", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return nil, err
	}

	return converter.WithdrawalToResponse(withdrawal), nil
}

func (u *withdrawalUseCase) FindAll(ctx context.Context) (*[]*model.WithdrawalResponse, error) {
	withdrawals, err := u.withdrawalRepository.FindAll(ctx, u.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid email or password")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to find user by email", err)
	}

	return converter.WithdrawalsToResponses(withdrawals), nil
}

func (u *withdrawalUseCase) FindById(ctx context.Context, request *model.FindWithdrawalById) (*model.WithdrawalResponse, error) {
	withdrawal, err := u.withdrawalRepository.FindById(ctx, u.db, request.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid email or password")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to find user by email", err)
	}

	return converter.WithdrawalToResponse(withdrawal), nil
}

func (u *withdrawalUseCase) Update(ctx context.Context, request *model.UpdateWithdrawalStatusRequest) (*model.WithdrawalResponse, error) {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return nil, err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	updateWitdrawal := &entity.Withdrawal{
		Id:     request.Id,
		Status: request.Status,
	}

	updateWitdrawal, err = u.withdrawalRepository.UpdateWithdrawalStatus(ctx, tx, updateWitdrawal)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to find update withdrawal status id", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return nil, err
	}

	return converter.WithdrawalToResponse(updateWitdrawal), nil
}

//cancel
//create
//retry
