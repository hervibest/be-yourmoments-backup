package usecase

import (
	"be-yourmoments/transaction-svc/internal/entity"
	errorcode "be-yourmoments/transaction-svc/internal/enum/error"
	"be-yourmoments/transaction-svc/internal/helper"
	"be-yourmoments/transaction-svc/internal/helper/logger"
	"be-yourmoments/transaction-svc/internal/helper/nullable"
	"be-yourmoments/transaction-svc/internal/model"
	"be-yourmoments/transaction-svc/internal/model/converter"
	"be-yourmoments/transaction-svc/internal/repository"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type ReviewUseCase interface {
	Create(ctx context.Context, request *model.CreateReviewRequest) (*model.CreatorReviewResponse, error)
	GetCreatorReview(ctx context.Context, request *model.GetAllReviewRequest) (*[]*model.CreatorReviewResponse, *model.PageMetadata, error)
}
type reviewUseCase struct {
	transactionDetailRepo repository.TransactionDetailRepository
	creatorReviewRepo     repository.CreatorReviewRepository
	db                    *sqlx.DB
	logs                  *logger.Log
}

func NewReviewUseCase(transactionDetailRepo repository.TransactionDetailRepository, creatorReviewRepo repository.CreatorReviewRepository,
	db *sqlx.DB, logs *logger.Log) ReviewUseCase {
	return &reviewUseCase{
		transactionDetailRepo: transactionDetailRepo,
		creatorReviewRepo:     creatorReviewRepo,
		db:                    db,
		logs:                  logs,
	}
}

func (u *reviewUseCase) Create(ctx context.Context, request *model.CreateReviewRequest) (*model.CreatorReviewResponse, error) {
	transactionDetail, err := u.transactionDetailRepo.FindByID(ctx, u.db, request.TransactionDetailId)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid transaction detail id")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to find transaction detail by transaction detail id", err)
	}

	if transactionDetail.IsReviewed {
		return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "User has reviewed the creator")
	}

	now := time.Now()
	review := &entity.CreatorReview{
		Id:                  ulid.Make().String(),
		TransactionDetailId: request.TransactionDetailId,
		CreatorId:           request.CreatorId,
		UserId:              request.UserId,
		Star:                request.Star,
		Comment:             nullable.ToSQLString(request.Comment),
		CreatedAt:           &now,
		UpdatedAt:           &now,
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return nil, err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	review, err = u.creatorReviewRepo.Create(ctx, tx, review)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to create creator review in database", err)
	}

	transactionDetail.IsReviewed = true
	_, err = u.transactionDetailRepo.UpdateReviewStatus(ctx, tx, transactionDetail)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to update transaction detail review status in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return nil, err
	}

	return converter.ReviewToResponse(review), err
}

func (u *reviewUseCase) GetCreatorReview(ctx context.Context, request *model.GetAllReviewRequest) (*[]*model.CreatorReviewResponse, *model.PageMetadata, error) {
	userPublicChat, pageMetadata, err := u.creatorReviewRepo.FindAll(ctx, u.db, request.Page, request.Size, request.Star, request.Order)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to find all creator review", err)
	}

	return converter.ReviewsToResponses(&userPublicChat), pageMetadata, nil
}
