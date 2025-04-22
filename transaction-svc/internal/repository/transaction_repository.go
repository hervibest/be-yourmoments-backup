package repository

import (
	"be-yourmoments/transaction-svc/internal/entity"
	"context"
	"fmt"
	"log"
)

type TransactionRepository interface {
	Create(ctx context.Context, db Querier, transaction *entity.Transaction) (*entity.Transaction, error)
	UpdateToken(ctx context.Context, db Querier, transaction *entity.Transaction) error
	FindById(ctx context.Context, db Querier, transactionId string) (*entity.Transaction, error)
	UpdateCallback(ctx context.Context, db Querier, transaction *entity.Transaction) error
}

type transactionRepository struct {
}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{}
}

func (r *transactionRepository) Create(ctx context.Context, db Querier, transaction *entity.Transaction) (*entity.Transaction, error) {
	query := `INSERT INTO transactions 
			  (id, user_id, status, photo_ids, checkout_at, amount, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := db.ExecContext(ctx, query, transaction.Id, transaction.UserId, transaction.Status, transaction.PhotoIds, transaction.CheckoutAt,
		transaction.Amount, transaction.CreatedAt, transaction.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert transactions: %w", err)
	}

	return transaction, nil
}

func (r *transactionRepository) UpdateToken(ctx context.Context, db Querier, transaction *entity.Transaction) error {
	query := `UPDATE transactions SET snap_token = $1, updated_at = $2 WHERE id = $3`

	_, err := db.ExecContext(ctx, query, transaction.SnapToken, transaction.UpdatedAt, transaction.Id)
	if err != nil {
		return fmt.Errorf("failed to update transaction token: %w", err)
	}

	return nil
}

func (r *transactionRepository) UpdateCallback(ctx context.Context, db Querier, transaction *entity.Transaction) error {
	query := `
	UPDATE transactions 
	SET 
		status = $1,
		payment_at = COALESCE($2, payment_at),
		snap_token = COALESCE($3, snap_token),
		external_status = COALESCE($4, external_status),
		external_callback_response = COALESCE($5, external_callback_response),
		updated_at = $6
	WHERE id = $7
	`
	_, err := db.ExecContext(ctx, query, transaction.Status, transaction.PaymentAt, transaction.SnapToken, transaction.ExternalStatus, transaction.ExternalCallbackResponse, transaction.UpdatedAt, transaction.Id)
	if err != nil {
		return fmt.Errorf("failed to update transaction callback: %w", err)
	}

	return nil
}

func (r *transactionRepository) FindById(ctx context.Context, db Querier, transactionId string) (*entity.Transaction, error) {
	transaction := new(entity.Transaction)
	query := `SELECT * FROM transactions WHERE id = $1`
	if err := db.GetContext(ctx, transaction, query, transactionId); err != nil {
		log.Printf("Error happen in FindById with error : %s", err.Error())
		return nil, err
	}

	return transaction, nil
}
