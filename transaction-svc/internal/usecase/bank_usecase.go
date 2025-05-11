package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	errorcode "github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/nullable"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model/converter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type BankUseCase interface {
	Create(ctx context.Context, request *model.CreateBankRequest) (*model.BankResponse, error)
	FindById(ctx context.Context, request *model.FindBankByIdRequest) (*model.BankResponse, error)
	FindAll(ctx context.Context) (*[]*model.BankResponse, error)
	Update(ctx context.Context) (*model.BankResponse, error)
	Delete(ctx context.Context, request *model.DeleteBankRequest) error
}

type bankUseCase struct {
	db             *sqlx.DB
	bankRepository repository.BankRepository
	logs           *logger.Log
}

func NewBankUseCase(db *sqlx.DB, bankRepository repository.BankRepository, logs *logger.Log) BankUseCase {
	return &bankUseCase{db: db, bankRepository: bankRepository, logs: logs}
}

func (u *bankUseCase) Create(ctx context.Context, request *model.CreateBankRequest) (*model.BankResponse, error) {
	now := time.Now()
	bank := &entity.Bank{
		Id:        ulid.Make().String(),
		BankCode:  request.BankCode,
		Name:      request.Name,
		Alias:     nullable.ToSQLString(request.Alias),
		SwiftCode: nullable.ToSQLString(request.SwiftCode),
		LogoUrl:   nullable.ToSQLString(request.LogoUrl),
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

	bank, err = u.bankRepository.Create(ctx, tx, bank)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to create wallet in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return nil, err
	}

	return converter.BankToResponse(bank), nil
}

func (u *bankUseCase) FindById(ctx context.Context, request *model.FindBankByIdRequest) (*model.BankResponse, error) {

	bank, err := u.bankRepository.FindById(ctx, u.db, request.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid Bank ID")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to find user by email", err)
	}

	return converter.BankToResponse(bank), err
}

func (u *bankUseCase) FindAll(ctx context.Context) (*[]*model.BankResponse, error) {
	banks, err := u.bankRepository.FindAll(ctx, u.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid Bank ID")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to find all bank in database", err)
	}

	return converter.BanksToResponses(banks), nil
}

func (u *bankUseCase) Update(ctx context.Context) (*model.BankResponse, error) {
	return nil, nil
}

func (u *bankUseCase) Delete(ctx context.Context, request *model.DeleteBankRequest) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	_, err = u.bankRepository.FindById(ctx, u.db, request.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewUseCaseError(errorcode.ErrResourceNotFound, "Invalid Bank ID")
		}
		return helper.WrapInternalServerError(u.logs, "failed to find user by email", err)
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if err = u.bankRepository.Delete(ctx, tx, request.Id); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to create wallet in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}
