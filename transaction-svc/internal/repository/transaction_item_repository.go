package repository

import (
	"context"
	"log"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
)

type TransactionItemRepository interface {
	Create(ctx context.Context, db Querier, items []*entity.TransactionItem) (*[]*entity.TransactionItem, error)
}

type transactionItemRepository struct {
}

func NewTransactionItemRepository() TransactionItemRepository {
	return &transactionItemRepository{}
}

func (r *transactionItemRepository) Create(ctx context.Context, db Querier, items []*entity.TransactionItem) (*[]*entity.TransactionItem, error) {
	query := `INSERT INTO transaction_items (id, transaction_detail_id, photo_id, price, discount, final_price, created_at, updated_at)
	          VALUES (:id, :transaction_detail_id, :photo_id, :price, :discount, :final_price, :created_at, :updated_at)`

	_, err := db.NamedExecContext(ctx, query, items)
	if err != nil {
		log.Printf("error inserting transaction items: %v", err)
		return nil, err
	}

	return &items, nil
}
