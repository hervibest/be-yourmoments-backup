package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
	errorcode "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/converter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/repository"

	"github.com/jmoiron/sqlx"
)

type CheckoutUseCase interface {
	PreviewCheckout(ctx context.Context, request *model.PreviewCheckoutRequest) (*model.PreviewCheckoutResponse, error)
	OwnerOwnPhotos(ctx context.Context, request *model.OwnerOwnPhotosRequest) error
	LockPhotosAndCalculatePrice(ctx context.Context, request *model.CalculateRequest) (*[]*model.CheckoutItem, *model.Total, error)
	CancelPhotos(ctx context.Context, request *model.CancelPhotosRequest) error
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

func (u *checkoutUseCase) PreviewCheckout(ctx context.Context, previewRequest *model.PreviewCheckoutRequest) (*model.PreviewCheckoutResponse, error) {
	now := time.Now()
	request := &model.CalculateRequest{
		UserId:   previewRequest.UserId,
		PhotoIds: previewRequest.PhotoIds,
	}

	result, total, err := u.calculatePrice(ctx, u.db, request, false)
	if err != nil {
		return nil, err
	}

	return converter.CheckoutItemToResponse(result, total.Price, total.Discount, &now), nil
}

func (u *checkoutUseCase) LockPhotosAndCalculatePrice(ctx context.Context, request *model.CalculateRequest) (*[]*model.CheckoutItem, *model.Total, error) {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return nil, nil, err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	result, total, err := u.calculatePrice(ctx, tx, request, true)
	if err != nil {
		return nil, nil, err
	}

	if err := u.photoRepository.UpdatePhotoStatusesByIDs(ctx, tx, enum.PhotoStatusInTransactionEnum, request.PhotoIds); err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to update photo statuses by photo ids with status IN_TRANSACTION ", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return nil, nil, err
	}

	return result, total, err
}

func (u *checkoutUseCase) calculatePrice(ctx context.Context, tx repository.Querier, request *model.CalculateRequest, isTransaction bool) (*[]*model.CheckoutItem, *model.Total, error) {
	photos, err := u.photoRepository.GetSimilarPhotosByIDs(ctx, tx, request.UserId, request.CreatorId, request.PhotoIds, isTransaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid photo id")
		}
		return nil, nil, helper.WrapInternalServerError(u.logs, "error get photos by ids", err)
	}

	if len(*photos) == 0 {
		return nil, nil, helper.NewUseCaseError(errorcode.ErrResourceNotFound, "No available photos found")
	}

	if len(*photos) != len(request.PhotoIds) {
		return nil, nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Some photos is missing, please try again later")
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

func (u *checkoutUseCase) OwnerOwnPhotos(ctx context.Context, request *model.OwnerOwnPhotosRequest) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	photos, err := u.photoRepository.GetSimilarPhotosByIDsWithoutCreatorFilter(ctx, tx, request.OwnerId, request.PhotoIds, true)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "error get photos by ids", err)
	}

	if len(*photos) == 0 {
		return helper.NewUseCaseError(errorcode.ErrResourceNotFound, "No available photos found")
	}

	if len(*photos) != len(request.PhotoIds) {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Some photos is missing, please try again later")
	}

	if err := u.photoRepository.UpdatePhotoOwnerAndStatusByIds(ctx, tx, request.OwnerId, request.PhotoIds); err != nil {
		return helper.WrapInternalServerError(u.logs, "error update photo owner and status by photo ids", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}

func (u *checkoutUseCase) CancelPhotos(ctx context.Context, request *model.CancelPhotosRequest) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	photos, err := u.photoRepository.GetSimilarPhotosByIDsWithoutCreatorFilter(ctx, tx, request.UserId, request.PhotoIds, true)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "error get photos by ids", err)
	}

	if len(*photos) == 0 {
		return helper.NewUseCaseError(errorcode.ErrResourceNotFound, "No available photos found")
	}

	if len(*photos) != len(request.PhotoIds) {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Some photos is missing, please try again later")
	}

	if err := u.photoRepository.UpdatePhotoStatusesByIDs(ctx, tx, enum.PhotoStatusAvailableEnum, request.PhotoIds); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update photo statuses by photo ids with status AVAILABLE ", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return err
}
