package usecase

import (
	"be-yourmoments/photo-svc/internal/entity"
	"be-yourmoments/photo-svc/internal/enum"
	errorcode "be-yourmoments/photo-svc/internal/enum/error"
	"be-yourmoments/photo-svc/internal/helper"
	"be-yourmoments/photo-svc/internal/helper/logger"
	"be-yourmoments/photo-svc/internal/model"
	"be-yourmoments/photo-svc/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type CheckoutUseCase interface {
	PreviewCheckout(ctx context.Context, request *model.PreviewCheckoutRequest) (*model.PreviewCheckoutResponse, error)
	CalculatePrice(ctx context.Context, request *model.PreviewCheckoutRequest) (*[]*model.CheckoutItem, *model.Total, error)
	OwnerOwnPhotos(ctx context.Context, ownerID string, photoIds []string) error
}

type checkoutUseCase struct {
	db                        *sqlx.DB
	photoRepository           repository.PhotoRepository
	creatorRepository         repository.CreatorRepository
	creatorDiscountRepository repository.CreatorDiscountRepository
	logs                      *logger.Log
}

func NewCheckoutUseCase(db *sqlx.DB, photoRepository repository.PhotoRepository, creatorRepository repository.CreatorRepository,
	creatorDiscountRepository repository.CreatorDiscountRepository, logs *logger.Log) CheckoutUseCase {
	return &checkoutUseCase{
		db:                        db,
		photoRepository:           photoRepository,
		creatorRepository:         creatorRepository,
		creatorDiscountRepository: creatorDiscountRepository,
		logs:                      logs,
	}
}

func (u *checkoutUseCase) PreviewCheckout(ctx context.Context, request *model.PreviewCheckoutRequest) (*model.PreviewCheckoutResponse, error) {
	now := time.Now()
	result, total, err := u.CalculatePrice(ctx, request)
	if err != nil {
		return nil, err
	}

	return &model.PreviewCheckoutResponse{
		Items:         result,
		TotalPrice:    total.Price,
		TotalDiscount: total.Discount,
		CreatedAt:     &now,
	}, nil
}

func (u *checkoutUseCase) CalculatePrice(ctx context.Context, request *model.PreviewCheckoutRequest) (*[]*model.CheckoutItem, *model.Total, error) {
	//ISSUE #3 creator_id should not checked (redudant from auth middleware)
	//Find creator to make sure creator cannot buy their own photos
	creator, err := u.creatorRepository.FindByUserId(ctx, request.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid user id")
		}
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to find creator by user id", err)
	}

	// TODO tambahkan permistic locking ? dengan db transaction
	log.Print("creator id", creator.Id, request.PhotoIds)
	photos, err := u.photoRepository.GetSimilarPhotosByIDs(ctx, request.UserId, creator.Id, request.PhotoIds)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid photo id")
		}
		return nil, nil, helper.WrapInternalServerError(u.logs, "error get photos by ids", err)
	}

	// 1. Calculate each creator's photos and save to map
	/* ex {"xxi" : 3, "xx2" : 4}  */
	photoCount := make(map[string]int)
	for _, p := range *photos {
		photoCount[p.CreatorId]++
	}

	if len(*photos) != len(request.PhotoIds) {
		foundPhotoMap := make(map[string]bool)
		countPhotoMap := make(map[string]int)
		for _, p := range *photos {
			foundPhotoMap[p.Id] = true
			countPhotoMap[p.Id]++
		}

		notFoundIds := []string{}
		doubleIds := []string{}

		for _, id := range request.PhotoIds {
			if !foundPhotoMap[id] {
				notFoundIds = append(notFoundIds, id)
			}
			countPhotoMap[id]++
			if countPhotoMap[id] > 2 {
				doubleIds = append(doubleIds, id)
			}
		}
		if len(doubleIds) != 0 {
			return nil, nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, fmt.Sprintf("Double photo id : %s", doubleIds))
		} else {
			return nil, nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, fmt.Sprintf("Invalid photo ids : %s", notFoundIds))
		}
	}

	// 2. Save unique creatorIds to Array
	var creatorIds []string
	for id := range photoCount {
		creatorIds = append(creatorIds, id)
	}

	discountRules, err := u.creatorDiscountRepository.GetDiscountRules(ctx, creatorIds)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid discount rule")
		}
		return nil, nil, helper.WrapInternalServerError(u.logs, "error get discount rules by ids", err)
	}

	//Get best choice for creator discount
	discountMap := make(map[string]*entity.CreatorDiscount)
	for _, rule := range *discountRules {
		if photoCount[rule.CreatorId] >= rule.MinQuantity {
			if _, exists := discountMap[rule.CreatorId]; !exists {
				discountMap[rule.CreatorId] = rule
			}
		}
	}

	result := make([]*model.CheckoutItem, 0)
	var totalAmount int32 = 0
	var totalDiscount int32 = 0

	for _, p := range *photos {
		if disc, ok := discountMap[p.CreatorId]; ok {
			var discount int32 = 0
			if disc.DiscountType == enum.DiscountTypeFlat {
				discount = disc.Value
			} else if disc.DiscountType == enum.DiscountTypePercent {
				discount = p.Price * disc.Value / 100
			}

			final := p.Price - discount
			item := &model.CheckoutItem{
				PhotoId:             p.Id,
				CreatorId:           p.CreatorId,
				Title:               p.Title,
				YourMomentsUrl:      p.YourMomentsUrl.String,
				Price:               p.Price,
				Discount:            discount,
				DiscountMinQuantity: disc.MinQuantity,
				DiscountValue:       disc.Value,
				DiscountId:          disc.Id,
				DiscountType:        disc.DiscountType,
				FinalPrice:          final,
			}

			result = append(result, item)
			totalAmount += final
			totalDiscount += discount
		} else {
			item := &model.CheckoutItem{
				PhotoId:        p.Id,
				CreatorId:      p.CreatorId,
				Title:          p.Title,
				YourMomentsUrl: p.YourMomentsUrl.String,
				Price:          p.Price,
				FinalPrice:     p.Price,
			}

			result = append(result, item)
			totalAmount += p.Price
		}
	}

	total := &model.Total{
		Price:    totalAmount,
		Discount: totalDiscount,
	}

	return &result, total, nil
}

func (u *checkoutUseCase) OwnerOwnPhotos(ctx context.Context, ownerID string, photoIds []string) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if err := u.photoRepository.UpdatePhotoOwnerByPhotoIds(ctx, tx, ownerID, photoIds); err != nil {
		return helper.WrapInternalServerError(u.logs, "error update photo owner by photo ids", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}
