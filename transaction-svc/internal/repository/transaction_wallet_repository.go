package repository

import (
	"be-yourmoments/transaction-svc/internal/entity"
	"be-yourmoments/transaction-svc/internal/helper"
	"be-yourmoments/transaction-svc/internal/model"
	"context"
	"fmt"
	"strconv"
	"strings"
)

type TransactionWalletRepository interface {
	BulkInsert(ctx context.Context, db Querier, entries []*entity.TransactionWallet) error
	FindAllByWalletId(ctx context.Context, tx Querier, page, size int, walletId, max, min, timeOrder string) ([]*entity.TransactionWallet, *model.PageMetadata, error)
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

func (r *transactionWalletRepository) FindManyByWalletId(ctx context.Context, db Querier, walletId string) (*[]*entity.TransactionWallet, error) {
	transactionWallets := make([]*entity.TransactionWallet, 0)
	const query = `
	SELECT * FROM transaction_wallets AS tw
	WHERE tw.wallet_id = $1
	`
	if err := db.SelectContext(ctx, transactionWallets, query, walletId); err != nil {
		return nil, fmt.Errorf("select context find many by wallet id error : %+v", err)
	}

	return &transactionWallets, nil
}

func (r *transactionWalletRepository) FindAllByWalletId(
	ctx context.Context,
	tx Querier,
	page, size int,
	walletId, max, min, timeOrder string,
) ([]*entity.TransactionWallet, *model.PageMetadata, error) {
	results := make([]*entity.TransactionWallet, 0)
	var totalItems int

	query := `SELECT * FROM transaction_wallets`
	countQuery := `SELECT COUNT(*) FROM transaction_wallets`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if walletId != "" {
		conditions = append(conditions, fmt.Sprintf("wallet_id = $%d", argIndex))
		args = append(args, walletId)
		argIndex++
	}

	if min != "" {
		if minVal, err := strconv.Atoi(min); err == nil {
			conditions = append(conditions, fmt.Sprintf("amount >= $%d", argIndex))
			args = append(args, minVal)
			argIndex++
		}
	}

	if max != "" {
		if maxVal, err := strconv.Atoi(max); err == nil {
			conditions = append(conditions, fmt.Sprintf("amount <= $%d", argIndex))
			args = append(args, maxVal)
			argIndex++
		}
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		query += whereClause
		countQuery += whereClause
	}

	order := "DESC"
	if strings.ToUpper(timeOrder) == "ASC" {
		order = "ASC"
	}
	query += fmt.Sprintf(" ORDER BY created_at %s", order)

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, size, (page-1)*size)

	if err := tx.GetContext(ctx, &totalItems, countQuery, args[:argIndex-1]...); err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(int64(totalItems), page, size)

	if err := tx.SelectContext(ctx, &results, query, args...); err != nil {
		return nil, nil, err
	}

	return results, pageMetadata, nil
}
