package repository

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/enum"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx Querier, transaction *entity.Transaction) (*entity.Transaction, error)
	UpdateToken(ctx context.Context, tx Querier, transaction *entity.Transaction) error
	FindById(ctx context.Context, tx Querier, transactionId string) (*entity.Transaction, error)
	UpdateCallback(ctx context.Context, tx Querier, transaction *entity.Transaction) error
	UserFindWithDetailById(ctx context.Context, tx Querier, transactionId, userId string) (*[]*entity.TransactionWithDetail, error)
	UserFindAll(ctx context.Context, tx Querier, page, size int, userId string, timeOrder string) (*[]*entity.Transaction, *model.PageMetadata, error)
	UpdateStatus(ctx context.Context, tx Querier, transaction *entity.Transaction) error
	FindManyByStatus(ctx context.Context, tx Querier, status enum.TransactionStatus) (*[]*entity.Transaction, error)
}

type transactionRepository struct {
}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{}
}

func (r *transactionRepository) Create(ctx context.Context, tx Querier, transaction *entity.Transaction) (*entity.Transaction, error) {
	query := `INSERT INTO transactions 
			  (id, user_id, status, photo_ids, checkout_at, amount, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := tx.ExecContext(ctx, query, transaction.Id, transaction.UserId, transaction.Status, transaction.PhotoIds, transaction.CheckoutAt,
		transaction.Amount, transaction.CreatedAt, transaction.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert transactions: %w", err)
	}

	return transaction, nil
}

func (r *transactionRepository) UpdateToken(ctx context.Context, tx Querier, transaction *entity.Transaction) error {
	query := `UPDATE transactions SET snap_token = $1, updated_at = $2 WHERE id = $3`

	_, err := tx.ExecContext(ctx, query, transaction.SnapToken, transaction.UpdatedAt, transaction.Id)
	if err != nil {
		return fmt.Errorf("failed to update transaction token: %w", err)
	}

	return nil
}

func (r *transactionRepository) UpdateCallback(ctx context.Context, tx Querier, transaction *entity.Transaction) error {
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
	_, err := tx.ExecContext(ctx, query, transaction.Status, transaction.PaymentAt, transaction.SnapToken, transaction.ExternalStatus, transaction.ExternalCallbackResponse, transaction.UpdatedAt, transaction.Id)
	if err != nil {
		return fmt.Errorf("failed to update transaction callback: %w", err)
	}

	return nil
}

func (r *transactionRepository) FindById(ctx context.Context, tx Querier, transactionId string) (*entity.Transaction, error) {
	transaction := new(entity.Transaction)
	query := `SELECT * FROM transactions WHERE id = $1`
	if err := tx.GetContext(ctx, transaction, query, transactionId); err != nil {
		log.Printf("Error happen in FindById with error : %s", err.Error())
		return nil, err
	}

	return transaction, nil
}

func (r *transactionRepository) UserFindWithDetailById(ctx context.Context, tx Querier, transactionId, userId string) (*[]*entity.TransactionWithDetail, error) {
	transactionWithDetails := make([]*entity.TransactionWithDetail, 0)
	query := `
	SELECT 
		trx.id AS transaction_id,
		trx.user_id,
		trx.status,
		trx.transaction_method_id,
		trx.transaction_type_id,
		trx.payment_type_id,
		trx.photo_ids,
		trx.payment_at,
		trx.checkout_at,
		trx.snap_token,
		-- trx.external_status,
		-- trx.external_callback_response,
		trx.amount,
		trx.created_at AS transaction_created_at,
		trx.updated_at AS transaction_updated_at,

		td.id AS transaction_detail_id,
		td.creator_id,
		td.creator_discount_id,
		td.is_reviewed,

		ti.id AS transaction_item_id,
		ti.photo_id,
		ti.price,
		ti.discount,
		ti.final_price
	FROM 
	 	transactions AS trx
	JOIN 
	 	transaction_details AS td
	ON
	 	trx.id = td.transaction_id
	JOIN
		transaction_items AS ti
	ON
	 	td.id = ti.transaction_detail_id
	WHERE 
		trx.id = $1
	AND
		trx.user_id = $2`
	if err := tx.SelectContext(ctx, &transactionWithDetails, query, transactionId, userId); err != nil {
		log.Printf("Error happen in FindById with error : %s", err.Error())
		return nil, err
	}

	return &transactionWithDetails, nil
}

func (r *transactionRepository) UserFindAll(ctx context.Context, tx Querier, page, size int, userId string, timeOrder string) (*[]*entity.Transaction, *model.PageMetadata, error) {
	results := make([]*entity.Transaction, 0)
	var totalItems int

	var conditions []string
	var args []interface{}
	argIndex := 1

	// Filter user_id
	conditions = append(conditions, "user_id = $"+strconv.Itoa(argIndex))
	args = append(args, userId)
	argIndex++

	// WHERE clause
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Base query
	baseQuery := `
	SELECT 
		id,
		user_id,
		status,
		transaction_method_id,
		transaction_type_id,
		payment_type_id,
		payment_at,
		checkout_at,
		amount,
		created_at,
		updated_at
	FROM 
		transactions` + whereClause

	// Count query
	countQuery := `SELECT COUNT(*) FROM transactions` + whereClause

	// Order by created_at
	if strings.ToUpper(timeOrder) == "ASC" || strings.ToUpper(timeOrder) == "DESC" {
		baseQuery += " ORDER BY created_at " + strings.ToUpper(timeOrder)
	} else {
		baseQuery += " ORDER BY created_at DESC"
	}

	// Pagination
	baseQuery += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	argsWithPagination := append([]interface{}{}, args...) // clone args
	argsWithPagination = append(argsWithPagination, size, (page-1)*size)

	// Get total count
	if err := tx.GetContext(ctx, &totalItems, countQuery, args...); err != nil {
		return nil, nil, err
	}

	// Get paginated data
	if err := tx.SelectContext(ctx, &results, baseQuery, argsWithPagination...); err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(int64(totalItems), page, size)

	return &results, pageMetadata, nil
}

func (r *transactionRepository) UpdateStatus(ctx context.Context, tx Querier, transaction *entity.Transaction) error {
	query := `UPDATE transactions SET status = $1, snap_token = $2, updated_at = $3 WHERE id = $4`

	_, err := tx.ExecContext(ctx, query, transaction.Status, transaction.SnapToken, transaction.UpdatedAt, transaction.Id)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	return nil
}

func (r *transactionRepository) FindManyByStatus(ctx context.Context, tx Querier, status enum.TransactionStatus) (*[]*entity.Transaction, error) {
	transactions := make([]*entity.Transaction, 0)
	query := `
	SELECT
		id, user_id, status
	FROM
		transactions
	WHERE
		status = $1
	`
	if err := tx.SelectContext(ctx, &transactions, query, status); err != nil {
		return nil, err
	}
	return &transactions, nil
}
