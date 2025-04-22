package usecase

import (
	"be-yourmoments/transaction-svc/internal/entity"
	errorcode "be-yourmoments/transaction-svc/internal/enum/error"
	"be-yourmoments/transaction-svc/internal/helper"
	"be-yourmoments/transaction-svc/internal/helper/logger"
	"be-yourmoments/transaction-svc/internal/model"
	"be-yourmoments/transaction-svc/internal/model/converter"
	"be-yourmoments/transaction-svc/internal/repository"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type BankWalletUseCase interface {
	Create(ctx context.Context, request *model.CreateBankWalletRequest) (*model.BankWalletResponse, error)
	Delete(ctx context.Context, request *model.DeleteBankWalletRequest) error
	FindAll(ctx context.Context) (*[]*model.BankWalletResponse, error)
	Update(ctx context.Context) (*model.BankResponse, error)
}

type bankWalletUseCase struct {
	db             *sqlx.DB
	bankWalletRepo repository.BankWalletRepository
	logs           *logger.Log
}

func NewBankWalletUseCase(db *sqlx.DB, bankWalletRepo repository.BankWalletRepository, logs *logger.Log) BankWalletUseCase {
	return &bankWalletUseCase{db: db, bankWalletRepo: bankWalletRepo, logs: logs}
}

func (u *bankWalletUseCase) Create(ctx context.Context, request *model.CreateBankWalletRequest) (*model.BankWalletResponse, error) {
	now := time.Now()
	bankWallet := &entity.BankWallet{
		Id:            ulid.Make().String(),
		BankId:        request.BankId,
		WalletId:      request.WalletId,
		FullName:      request.FullName,
		AccountNumber: request.AccountNumber,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return nil, err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	bankWallet, err = u.bankWalletRepo.Create(ctx, tx, bankWallet)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to create wallet in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return nil, err
	}

	return converter.BankWalletToResponse(bankWallet), nil
}

// func (u *bankWalletUseCase) FindById(ctx context.Context, request *model.FindByIdRequest) (*model.BankResponse, error) {

// 	bank, err := u.bankWalletRepo.FindById(ctx, u.db, request.Id)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid email or password")
// 		}
// 		return nil, helper.WrapInternalServerError(u.logs, "failed to find user by email", err)
// 	}

// 	return converter.BankToResponse(bank), err
// }

func (u *bankWalletUseCase) FindAll(ctx context.Context) (*[]*model.BankWalletResponse, error) {
	bankWallets, err := u.bankWalletRepo.FindAll(ctx, u.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "invalid request")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to find all bank wallet", err)
	}

	return converter.BankWalletsToResponses(bankWallets), nil
}

func (u *bankWalletUseCase) Update(ctx context.Context) (*model.BankResponse, error) {
	return nil, nil
}

func (u *bankWalletUseCase) Delete(ctx context.Context, request *model.DeleteBankWalletRequest) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if err = u.bankWalletRepo.Delete(ctx, tx, request.Id); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to create wallet in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}
