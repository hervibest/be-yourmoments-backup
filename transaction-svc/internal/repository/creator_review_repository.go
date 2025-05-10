package repository

import (
	"context"
	"strconv"
	"strings"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
)

type CreatorReviewRepository interface {
	Create(ctx context.Context, tx Querier, review *entity.CreatorReview) (*entity.CreatorReview, error)
	FindAll(ctx context.Context, tx Querier, page int, size int, rating int, timeOrder string) ([]*entity.CreatorReview, *model.PageMetadata, error)
	CountTotalReviewAndRating(ctx context.Context, tx Querier, creatorId string) (*entity.TotalReviewAndRating, error)
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
	rating, 
	comment, 
	created_at, 
	updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	`
	_, err := tx.ExecContext(ctx, query, review.Id, review.TransactionDetailId, review.CreatorId, review.UserId, review.Rating, review.Comment, review.CreatedAt, review.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return review, nil
}

// TODO tambahkan get by creator id atau kebalikkanya
func (r *creatorReviewRepository) FindAll(ctx context.Context, tx Querier, page, size int, rating int, timeOrder string) ([]*entity.CreatorReview, *model.PageMetadata, error) {
	results := make([]*entity.CreatorReview, 0)
	var totalItems int

	query := `SELECT * FROM creator_reviews`
	countQuery := `SELECT COUNT(*) FROM creator_reviews`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if rating != 0 {
		conditions = append(conditions, "rating = $"+strconv.Itoa(argIndex))
		args = append(args, rating)
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

	if err := tx.SelectContext(ctx, &results, query, args...); err != nil {
		return nil, nil, err
	}

	return results, pageMetadata, nil
}

func (r *creatorReviewRepository) CountTotalReviewAndRating(ctx context.Context, tx Querier, creatorId string) (*entity.TotalReviewAndRating, error) {
	totalReviewAndRating := new(entity.TotalReviewAndRating)

	query := `
	SELECT
		COUNT(*) AS total_review, 
		AVG(rating) AS average_rating
	FROM 
		creator_reviews 
	WHERE 
		creator_id = $1
	`

	if err := tx.GetContext(ctx, totalReviewAndRating, query, creatorId); err != nil {
		return nil, err
	}

	return totalReviewAndRating, nil
}
