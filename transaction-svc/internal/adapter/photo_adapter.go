package adapter

import (
	"context"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/discovery"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/utils"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
)

type PhotoAdapter interface {
	CalculatePhotoPrice(ctx context.Context, userId, creatorId string, photoIds []string) (*[]*model.CheckoutItem, *model.Total, error)
	OwnerOwnPhotos(ctx context.Context, ownerId string, photoIds []string) error
	GetPhotoWithDetails(ctx context.Context, photoIds []string, userId string) (*[]*photopb.Photo, error)
	CancelPhotos(ctx context.Context, userId string, photoIds []string) error
}

type photoAdapter struct {
	client photopb.PhotoServiceClient
}

func NewPhotoAdapter(ctx context.Context, registry discovery.Registry, logs *logger.Log) (PhotoAdapter, error) {
	photoServiceName := utils.GetEnv("PHOTO_SVC_NAME")
	conn, err := discovery.ServiceConnection(ctx, photoServiceName, registry, logs)
	if err != nil {
		return nil, err
	}

	client := photopb.NewPhotoServiceClient(conn)

	return &photoAdapter{
		client: client,
	}, nil
}

func (a *photoAdapter) CalculatePhotoPrice(ctx context.Context, userId, creatorId string, photoIds []string) (*[]*model.CheckoutItem, *model.Total, error) {
	processPhotoRequest := &photopb.CalculatePhotoPriceRequest{
		UserId:    userId,
		CreatorId: creatorId,
		PhotoIds:  photoIds,
	}

	response, err := a.client.CalculatePhotoPrice(ctx, processPhotoRequest)
	if err != nil {
		return nil, nil, helper.FromGRPCError(err)
	}

	items := make([]*model.CheckoutItem, 0)
	for _, item := range response.Items {
		transactionItem := &model.CheckoutItem{
			PhotoId:             item.GetPhotoId(),
			CreatorId:           item.GetCreatorId(),
			Title:               item.GetTitle(),
			YourMomentsUrl:      item.GetYourMomentsUrl(),
			Price:               item.GetPrice(),
			Discount:            item.GetDiscount(),
			DiscountMinQuantity: int(item.GetDiscountMinQuantity()),
			DiscountValue:       item.GetDiscountValue(),
			DiscountId:          item.GetDiscountId(),
			DiscountType:        item.GetDiscountType(),
			FinalPrice:          item.GetFinalPrice(),
		}
		items = append(items, transactionItem)
	}

	total := &model.Total{
		Price:    response.Total.GetPrice(),
		Discount: response.Total.GetDiscount(),
	}

	return &items, total, nil
}

func (a *photoAdapter) OwnerOwnPhotos(ctx context.Context, ownerId string, photoIds []string) error {
	ownerOwnPhotosRequest := &photopb.OwnerOwnPhotosRequest{
		OwnerId:  ownerId,
		PhotoIds: photoIds,
	}

	_, err := a.client.OwnerOwnPhotos(ctx, ownerOwnPhotosRequest)
	if err != nil {
		return helper.FromGRPCError(err)
	}

	return nil
}

func (a *photoAdapter) CancelPhotos(ctx context.Context, userId string, photoIds []string) error {
	cancelPhotosRequest := &photopb.CancelPhotosRequest{
		UserId:   userId,
		PhotoIds: photoIds,
	}

	_, err := a.client.CancelPhotos(ctx, cancelPhotosRequest)
	if err != nil {
		return helper.FromGRPCError(err)
	}

	return nil
}

func (a *photoAdapter) GetPhotoWithDetails(ctx context.Context, photoIds []string, userId string) (*[]*photopb.Photo, error) {
	processPhotoRequest := &photopb.GetPhotoWithDetailsRequest{
		PhotoIds: photoIds,
		UserId:   userId,
	}

	response, err := a.client.GetPhotoWithDetails(ctx, processPhotoRequest)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	return &response.PhotoWithDetails, nil
}
