package repository

import (
	"be-yourmoments/transaction-svc/internal/entity"
	"context"
	"fmt"
	"strings"
)

type TransactionWalletRepository interface {
	BulkInsert(ctx context.Context, db Querier, entries []*entity.TransactionWallet) error
}
type transactionWalletRepository struct{}

func NewTransactionWalletRepository() TransactionWalletRepository {
	return &transactionWalletRepository{}
}

func (r *transactionWalletRepository) BulkInsert(ctx context.Context, db Querier, entries []*entity.TransactionWallet) error {
	if len(entries) == 0 {
		return nil
	}

	const query = `
        INSERT INTO transaction_wallets
        (id, wallet_id, transaction_detail_id, amount, created_at, updated_at)
        VALUES `

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for i, entry := range entries {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))
		valueArgs = append(valueArgs,
			entry.Id,
			entry.WalletId,
			entry.TransactionDetailId,
			entry.Amount,
			entry.CreatedAt,
			entry.UpdatedAt,
		)
	}

	fullQuery := query + strings.Join(valueStrings, ",")
	_, err := db.ExecContext(ctx, fullQuery, valueArgs...)
	return err
}
