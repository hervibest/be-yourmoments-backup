package repository

import (
	"context"
	"log"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
)

type TransactionDetailRepository interface {
	Create(ctx context.Context, tx Querier, trxId string, details []*entity.TransactionDetail) (*[]*entity.TransactionDetail, error)
	FindManyByTrxID(ctx context.Context, tx Querier, trxId string) (*[]*entity.TransactionDetail, error)
	FindByIDNotCreator(ctx context.Context, tx Querier, transactionDetailID, creatorID string, isForUpdate bool) (*entity.TransactionDetail, error)
	UpdateReviewStatus(ctx context.Context, tx Querier, transactionDetail *entity.TransactionDetail) (*entity.TransactionDetail, error)
}

type transactionDetailRepository struct {
}

func NewTransactionDetailRepository() TransactionDetailRepository {
	return &transactionDetailRepository{}
}

func (r *transactionDetailRepository) Create(ctx context.Context, tx Querier, trxId string, details []*entity.TransactionDetail) (*[]*entity.TransactionDetail, error) {
	query := `INSERT INTO transaction_details (id, transaction_id, creator_id, subtotal_price, creator_discount_id, created_at, updated_at)
	          VALUES (:id, :transaction_id, :creator_id, :subtotal_price, :creator_discount_id, :created_at, :updated_at)`
	for i := range details {
		details[i].TransactionId = trxId
	}

	_, err := tx.NamedExecContext(ctx, query, details)
	if err != nil {
		return nil, err
	}

	return &details, nil
}

func (r *transactionDetailRepository) FindManyByTrxID(ctx context.Context, tx Querier, trxId string) (*[]*entity.TransactionDetail, error) {
	const query = `
        SELECT id, transaction_id, creator_id, subtotal_price
        FROM transaction_details
        WHERE transaction_id = $1
    `
	details := make([]*entity.TransactionDetail, 0)
	err := tx.SelectContext(ctx, &details, query, trxId)
	if err != nil {
		return nil, err
	}

	return &details, nil
}

func (r *transactionDetailRepository) FindByIDNotCreator(ctx context.Context, tx Querier, transactionDetailID, creatorID string, isForUpdate bool) (*entity.TransactionDetail, error) {
	log.Default().Println("Creator ID in repo:", creatorID)

	transactionDetail := new(entity.TransactionDetail)
	var query = `
        SELECT id, transaction_id, creator_id, subtotal_price, is_reviewed
        FROM transaction_details
        WHERE id = $1 AND creator_id != $2
    `

	if isForUpdate {
		query += " FOR UPDATE"
	}

	err := tx.GetContext(ctx, transactionDetail, query, transactionDetailID, creatorID)
	if err != nil {
		return nil, err
	}

	return transactionDetail, nil
}

func (r *transactionDetailRepository) UpdateReviewStatus(ctx context.Context, tx Querier, transactionDetail *entity.TransactionDetail) (*entity.TransactionDetail, error) {
	const query = `
        UPDATE transaction_details SET is_reviewed = $1 WHERE id = $2
    `
	_, err := tx.ExecContext(ctx, query, transactionDetail.IsReviewed, transactionDetail.Id)
	if err != nil {
		return nil, err
	}

	return transactionDetail, nil
}
