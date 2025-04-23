package repository

import (
	"be-yourmoments/transaction-svc/internal/entity"
	"be-yourmoments/transaction-svc/internal/helper"
	"be-yourmoments/transaction-svc/internal/model"
	"context"
	"strconv"
	"strings"
)

type CreatorReviewRepository interface {
	Create(ctx context.Context, tx Querier, review *entity.CreatorReview) (*entity.CreatorReview, error)
	FindAll(ctx context.Context, tx Querier, page int, size int, star int, timeOrder string) ([]*entity.CreatorReview, *model.PageMetadata, error)
}
type creatorReviewRepository struct{}

func NewCreatorReviewRepository() CreatorReviewRepository {
	return &creatorReviewRepository{}
}

func (r *creatorReviewRepository) Create(ctx context.Context, tx Querier, review *entity.CreatorReview) (*entity.CreatorReview, error) {
	query := `
	INSERT INTO creator_reviews 
	(id, transaction_detail_id, 
	creator_id, 
	user_id,
	star, 
	comment, 
	created_at, 
	updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	`
	_, err := tx.ExecContext(ctx, query, review.Id, review.TransactionDetailId, review.CreatorId, review.UserId, review.Star, review.Comment, review.CreatedAt, review.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return review, nil
}

// TODO tambahkan get by creator id atau kebalikkanya
func (r *creatorReviewRepository) FindAll(ctx context.Context, tx Querier, page, size int, star int, timeOrder string) ([]*entity.CreatorReview, *model.PageMetadata, error) {
	results := make([]*entity.CreatorReview, 0)
	var totalItems int

	query := `SELECT * FROM creator_reviews`
	countQuery := `SELECT COUNT(*) FROM creator_reviews`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if star != 0 {
		conditions = append(conditions, "star = $"+strconv.Itoa(argIndex))
		args = append(args, star)
		argIndex++
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		query += whereClause
		countQuery += whereClause
	}

	// ORDER BY created_at ASC or DESC
	if strings.ToUpper(timeOrder) == "ASC" || strings.ToUpper(timeOrder) == "DESC" {
		query += " ORDER BY created_at " + strings.ToUpper(timeOrder)
	} else {
		query += " ORDER BY created_at DESC" // default DESC if not specified
	}

	// Pagination
	query += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, size, (page-1)*size)

	// Total count
	if err := tx.GetContext(ctx, &totalItems, countQuery, args[:argIndex-1]...); err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(int64(totalItems), page, size)

	if err := tx.SelectContext(ctx, results, query, args...); err != nil {
		return nil, nil, err
	}

	return results, pageMetadata, nil
}
