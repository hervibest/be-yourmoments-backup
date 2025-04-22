package repository

import (
	"be-yourmoments/transaction-svc/internal/entity"
	"context"
)

type BankWalletRepository interface {
	Create(ctx context.Context, db Querier, bankWallet *entity.BankWallet) (*entity.BankWallet, error)
	Update(ctx context.Context, db Querier, bankWallet *entity.BankWallet) (*entity.BankWallet, error)
	// FindById(ctx context.Context, db Querier, bankId string) (*entity.BankWallet, error)
	FindAll(ctx context.Context, db Querier) (*[]*entity.BankWallet, error)
	Delete(ctx context.Context, db Querier, id string) error
}

type bankWalletRepository struct {
}

func NewBankWalletRepository() BankWalletRepository {
	return &bankWalletRepository{}
}

func (r *bankWalletRepository) Create(ctx context.Context, db Querier, bankWallet *entity.BankWallet) (*entity.BankWallet, error) {
	query := "INSERT INTO bank_wallets (id, wallet_id, bank_id, full_name, account_number, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := db.ExecContext(ctx, query, bankWallet.Id, bankWallet.WalletId, bankWallet.BankId, bankWallet.FullName, bankWallet.AccountNumber, bankWallet.CreatedAt, bankWallet.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return bankWallet, nil
}

func (r *bankWalletRepository) Update(ctx context.Context, db Querier, bankWallet *entity.BankWallet) (*entity.BankWallet, error) {
	return nil, nil
}

// func (r *bankWalletRepository) FindById(ctx context.Context, db Querier, bankId string) (*entity.BankWallet, error) {
// 	bank := new(entity.BankWallet)
// 	query := "SELECT * FROM bank_wallets WHERE id = $1"
// 	if err := db.GetContext(ctx, bank, query, bankId); err != nil {
// 		return nil, err
// 	}

// 	return bank, nil
// }

func (r *bankWalletRepository) FindAll(ctx context.Context, db Querier) (*[]*entity.BankWallet, error) {
	banks := make([]*entity.BankWallet, 0)
	query := "SELECT * FROM bank_wallets"
	if err := db.SelectContext(ctx, &banks, query); err != nil {
		return nil, err
	}

	return &banks, nil
}

func (r *bankWalletRepository) Delete(ctx context.Context, db Querier, id string) error {
	query := "DELETE FROM bank_wallets WHERE  id = $1"
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
