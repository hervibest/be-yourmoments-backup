package repository

import (
	"context"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
)

type BankRepository interface {
	Create(ctx context.Context, db Querier, bank *entity.Bank) (*entity.Bank, error)
	Update(ctx context.Context, db Querier, bank *entity.Bank) (*entity.Bank, error)
	FindById(ctx context.Context, db Querier, bankId string) (*entity.Bank, error)
	FindAll(ctx context.Context, db Querier) (*[]*entity.Bank, error)
	FindByCode(ctx context.Context, db Querier, bankCode string) (*entity.Bank, error)
	Delete(ctx context.Context, db Querier, bankId string) error
}

type bankRepository struct {
}

func NewBankRepository() BankRepository {
	return &bankRepository{}
}

func (r *bankRepository) Create(ctx context.Context, db Querier, bank *entity.Bank) (*entity.Bank, error) {
	query := "INSERT INTO banks (id, bank_code, name, alias, swift_code, logo_url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := db.ExecContext(ctx, query, bank.Id, bank.BankCode, bank.Name, bank.Alias, bank.SwiftCode, bank.LogoUrl, bank.CreatedAt, bank.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return bank, nil
}

func (r *bankRepository) Update(ctx context.Context, db Querier, bank *entity.Bank) (*entity.Bank, error) {
	return nil, nil
}

func (r *bankRepository) FindById(ctx context.Context, db Querier, bankId string) (*entity.Bank, error) {
	bank := new(entity.Bank)
	query := "SELECT * FROM banks WHERE id = $1"
	if err := db.GetContext(ctx, bank, query, bankId); err != nil {
		return nil, err
	}

	return bank, nil
}

func (r *bankRepository) FindAll(ctx context.Context, db Querier) (*[]*entity.Bank, error) {
	banks := make([]*entity.Bank, 0)
	query := "SELECT * FROM banks"
	if err := db.SelectContext(ctx, &banks, query); err != nil {
		return nil, err
	}

	return &banks, nil
}

func (r *bankRepository) Delete(ctx context.Context, db Querier, bankId string) error {
	query := "DELETE FROM banks WHERE id = $1"
	_, err := db.ExecContext(ctx, query, bankId)
	if err != nil {
		return err
	}
	return nil
}

func (r *bankRepository) FindByCode(ctx context.Context, db Querier, bankCode string) (*entity.Bank, error) {
	bank := new(entity.Bank)
	query := "SELECT * FROM banks WHERE bank_code = $1"
	if err := db.GetContext(ctx, bank, query, bankCode); err != nil {
		return nil, err
	}

	return bank, nil
}
