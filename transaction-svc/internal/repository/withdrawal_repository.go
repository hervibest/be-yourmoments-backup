package repository

import (
	"be-yourmoments/transaction-svc/internal/entity"
	"context"
)

type WithdrawalRepository interface {
	Create(ctx context.Context, db Querier, withdrawal *entity.Withdrawal) (*entity.Withdrawal, error)
	Delete(ctx context.Context, db Querier, withdrawalsId string) error
	FindAll(ctx context.Context, db Querier) (*[]*entity.Withdrawal, error)
	FindById(ctx context.Context, db Querier, bankId string) (*entity.Withdrawal, error)
	UpdateWithdrawalStatus(ctx context.Context, db Querier, withdrawal *entity.Withdrawal) (*entity.Withdrawal, error)
}

type withdrawalRepository struct {
}

func NewWithdrawalRepository() WithdrawalRepository {
	return &withdrawalRepository{}
}

func (r *withdrawalRepository) Create(ctx context.Context, db Querier, withdrawal *entity.Withdrawal) (*entity.Withdrawal, error) {
	query := "INSERT INTO withdrawals (id, bank_wallet_id, wallet_id, amount, status, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := db.ExecContext(ctx, query, withdrawal.Id, withdrawal.BankWalletId, withdrawal.WalletId, withdrawal.Amount, withdrawal.Status, withdrawal.Description, withdrawal.CreatedAt, withdrawal.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return withdrawal, nil
}

func (r *withdrawalRepository) UpdateWithdrawalStatus(ctx context.Context, db Querier, withdrawal *entity.Withdrawal) (*entity.Withdrawal, error) {
	query := "UPDATE withdrawals SET status = $1, description =  COALESCE($2, description), updated_at = $3 WHERE id = $4 RETURNING *"
	if err := db.GetContext(ctx, withdrawal, query, withdrawal.Status, withdrawal.Description, withdrawal.UpdatedAt, withdrawal.Id); err != nil {
		return nil, err
	}

	return withdrawal, nil
}

func (r *withdrawalRepository) FindById(ctx context.Context, db Querier, bankId string) (*entity.Withdrawal, error) {
	bank := new(entity.Withdrawal)
	query := "SELECT * FROM withdrawals WHERE id = $1"
	if err := db.GetContext(ctx, bank, query, bankId); err != nil {
		return nil, err
	}

	return bank, nil
}

func (r *withdrawalRepository) FindAll(ctx context.Context, db Querier) (*[]*entity.Withdrawal, error) {
	withdrawals := make([]*entity.Withdrawal, 0)
	query := "SELECT * FROM withdrawals"
	if err := db.SelectContext(ctx, withdrawals, query); err != nil {
		return nil, err
	}

	return &withdrawals, nil
}

func (r *withdrawalRepository) Delete(ctx context.Context, db Querier, withdrawalsId string) error {
	query := "DELETE FROM withdrawals WHERE id = $1"
	_, err := db.ExecContext(ctx, query, withdrawalsId)
	if err != nil {
		return err
	}
	return nil
}
