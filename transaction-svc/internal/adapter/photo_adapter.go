package adapter

import (
	"be-yourmoments/transaction-svc/internal/helper"
	"be-yourmoments/transaction-svc/internal/helper/discovery"
	"be-yourmoments/transaction-svc/internal/model"
	"context"

	photopb "github.com/be-yourmoments/pb/photo"
)

type PhotoAdapter interface {
	CalculatePhotoPrice(ctx context.Context, userId string, photoIds []string) (*[]*model.CheckoutItem, *model.Total, error)
	OwnerOwnPhotos(ctx context.Context, ownerId string, photoIds []string) error
}

type photoAdapter struct {
	client photopb.PhotoServiceClient
}

func NewPhotoAdapter(ctx context.Context, registry discovery.Registry) (PhotoAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "photo-svc-grpc", registry)
	if err != nil {
		return nil, err
	}

	client := photopb.NewPhotoServiceClient(conn)

	return &photoAdapter{
		client: client,
	}, nil
}

func (a *photoAdapter) CalculatePhotoPrice(ctx context.Context, userId string, photoIds []string) (*[]*model.CheckoutItem, *model.Total, error) {
	processPhotoRequest := &photopb.CalculatePhotoPriceRequest{
		UserId:   userId,
		PhotoIds: photoIds,
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
