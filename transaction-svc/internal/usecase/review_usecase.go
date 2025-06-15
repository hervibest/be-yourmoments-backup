package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	errorcode "github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/enum/error"
	producer "github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/gateway/messaging"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/nullable"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model/converter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model/event"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type ReviewUseCase interface {
	Create(ctx context.Context, request *model.CreateReviewRequest) (*model.CreatorReviewResponse, error)
	CreatorGetReview(ctx context.Context, request *model.GetAllReviewRequest) (*[]*model.CreatorReviewResponse, *model.PageMetadata, error)
}
type reviewUseCase struct {
	transactionDetailRepo repository.TransactionDetailRepository
	creatorReviewRepo     repository.CreatorReviewRepository
	transactionProducer   producer.TransactionProducer
	db                    *sqlx.DB
	logs                  *logger.Log
}

func NewReviewUseCase(transactionDetailRepo repository.TransactionDetailRepository, creatorReviewRepo repository.CreatorReviewRepository,
	transactionProducer producer.TransactionProducer, db *sqlx.DB, logs *logger.Log) ReviewUseCase {
	return &reviewUseCase{
		transactionDetailRepo: transactionDetailRepo,
		creatorReviewRepo:     creatorReviewRepo,
		transactionProducer:   transactionProducer,
		db:                    db,
		logs:                  logs,
	}
}

func (u *reviewUseCase) Create(ctx context.Context, request *model.CreateReviewRequest) (*model.CreatorReviewResponse, error) {
	transactionDetail, err := u.transactionDetailRepo.FindByID(ctx, u.db, request.TransactionDetailId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
		Rating:              request.Rating,
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

	totalReviewAndRating, err := u.creatorReviewRepo.CountTotalReviewAndRating(ctx, u.db, request.CreatorId)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to count total review and rating in database", err)
	}

	totalReviewAndRating.CreatorId = request.CreatorId

	creatorReviewCountEvent := &event.CreatorReviewCountEvent{
		Id:          request.CreatorId,
		Rating:      totalReviewAndRating.Rating,
		RatingCount: totalReviewAndRating.TotalReview,
	}

	if err := u.transactionProducer.ProduceCreateReviewEvent(ctx, creatorReviewCountEvent); err != nil {
		return nil, err
	}

	return converter.ReviewToResponse(review), err
}

func (u *reviewUseCase) CreatorGetReview(ctx context.Context, request *model.GetAllReviewRequest) (*[]*model.CreatorReviewResponse, *model.PageMetadata, error) {
	userPublicChat, pageMetadata, err := u.creatorReviewRepo.FindAll(ctx, u.db, request.Page, request.Size, request.Rating, request.CreatorId, request.Order)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to find all creator review", err)
	}

	return converter.ReviewsToResponses(&userPublicChat), pageMetadata, nil
}
