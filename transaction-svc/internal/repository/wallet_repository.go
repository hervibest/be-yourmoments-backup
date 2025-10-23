package repository

import (
	"context"
	"fmt"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"

	"github.com/lib/pq"
)

type WalletRepository interface {
	AddBalance(ctx context.Context, db Querier, walletID string, amount int64) error
	Create(ctx context.Context, db Querier, wallet *entity.Wallet) (*entity.Wallet, error)
	FindByCreatorIDs(ctx context.Context, db Querier, creatorIDs []string) (*[]*entity.Wallet, error)
	FindById(ctx context.Context, db Querier, walletId string, forUpdate bool) (*entity.Wallet, error)
	FindByCreatorId(ctx context.Context, db Querier, creatorId string) (*entity.Wallet, error)
	FindIdByCreatorId(ctx context.Context, db Querier, creatorId string) (string, error)
	ReduceBalance(ctx context.Context, db Querier, walletID string, amount int64) error
}

type walletRepository struct {
}

func NewWalletRepository() WalletRepository {
	return &walletRepository{}
}

func (r *walletRepository) Create(ctx context.Context, db Querier, wallet *entity.Wallet) (*entity.Wallet, error) {
	query := `
		INSERT INTO wallets 
		(id, creator_id, balance, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := db.ExecContext(ctx, query, wallet.Id, wallet.CreatorId, wallet.Balance, wallet.CreatedAt, wallet.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert wallet: %w", err)
	}

	return wallet, nil
}

func (r *walletRepository) FindByCreatorIDs(ctx context.Context, db Querier, creatorIDs []string) (*[]*entity.Wallet, error) {
	const query = `
        SELECT id, creator_id, balance
        FROM wallets
        WHERE creator_id = ANY($1)
    `
	wallets := make([]*entity.Wallet, 0)
	if err := db.SelectContext(ctx, &wallets, query, pq.Array(creatorIDs)); err != nil {
		return nil, err
	}

	return &wallets, nil
}

func (r *walletRepository) FindById(ctx context.Context, db Querier, walletId string, forUpdate bool) (*entity.Wallet, error) {
	query := `
        SELECT id, creator_id, balance, created_at, updated_at
        FROM wallets
        WHERE id = $1
    `
	if forUpdate {
		query += " FOR UDPATE"
	}
	wallet := new(entity.Wallet)
	err := db.GetContext(ctx, wallet, query, walletId)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (r *walletRepository) FindByCreatorId(ctx context.Context, db Querier, creatorId string) (*entity.Wallet, error) {
	const query = `
        SELECT id, creator_id, balance, created_at, updated_at
        FROM wallets
        WHERE creator_id = $1
    `
	wallet := new(entity.Wallet)
	err := db.GetContext(ctx, wallet, query, creatorId)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (r *walletRepository) AddBalance(ctx context.Context, db Querier, walletID string, amount int64) error {
	const query = `
        UPDATE wallets
        SET balance = balance + $1,
            updated_at = NOW()
        WHERE id = $2
    `
	_, err := db.ExecContext(ctx, query, amount, walletID)
	if err != nil {
		return err
	}
	return nil
}

func (r *walletRepository) ReduceBalance(ctx context.Context, db Querier, walletID string, amount int64) error {
	const query = `
        UPDATE wallets
        SET balance = balance - $1,
            updated_at = NOW()
        WHERE id = $2
    `
	_, err := db.ExecContext(ctx, query, amount, walletID)
	if err != nil {
		return err
	}
	return nil
}

func (r *walletRepository) FindIdByCreatorId(ctx context.Context, db Querier, creatorId string) (string, error) {
	const query = `
        SELECT id
        FROM wallets
        WHERE creator_id = $1
    `	
	var walletId string
	err := db.GetContext(ctx, &walletId, query, creatorId)
	if err != nil {
		return "", err
	}
	return walletId, nil
}
