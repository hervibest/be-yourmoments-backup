package adapter

import (
	"context"
	"fmt"
	"sync"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/discovery"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/utils"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"google.golang.org/grpc"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
)

type PhotoAdapter interface {
	CalculatePhotoPrice(ctx context.Context, userId, creatorId string, photoIds []string) (*[]*model.CheckoutItem, *model.Total, error)
	CalculatePhotoPriceV2(ctx context.Context, userId, creatorId string, request *model.CreateTransactionV2Request) (*[]*model.CheckoutItem, *model.Total, error)
	OwnerOwnPhotos(ctx context.Context, ownerId string, photoIds []string) error
	GetPhotoWithDetails(ctx context.Context, photoIds []string, userId string) (*[]*photopb.Photo, error)
	CancelPhotos(ctx context.Context, userId string, photoIds []string) error
	GetCreator(ctx context.Context, userId string) (*entity.Creator, error)
}

type photoAdapter struct {
	client   photopb.PhotoServiceClient
	conn     *grpc.ClientConn
	registry discovery.Registry
	logs     *logger.Log
	mu       sync.Mutex
}

func NewPhotoAdapter(ctx context.Context, registry discovery.Registry, logs *logger.Log) (PhotoAdapter, error) {
	photoServiceName := utils.GetEnv("PHOTO_SVC_NAME")
	conn, err := discovery.ServiceConnection(ctx, photoServiceName, registry, logs)
	if err != nil {
		return nil, err
	}

	client := photopb.NewPhotoServiceClient(conn)

	return &photoAdapter{
		client:   client,
		conn:     conn,
		registry: registry,
		logs:     logs,
	}, nil
}

// reconnect ke service photo via Consul
func (a *photoAdapter) reconnect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.logs.Log("reconnecting to photo service...")

	// Tutup koneksi lama agar tidak bocor
	if a.conn != nil {
		_ = a.conn.Close()
	}

	conn, err := discovery.ServiceConnection(ctx, utils.GetEnv("PHOTO_SVC_NAME"), a.registry, a.logs)
	if err != nil {
		a.logs.Error(fmt.Sprintf("reconnect failed: %v", err))
		return err
	}

	a.conn = conn
	a.client = photopb.NewPhotoServiceClient(conn)
	a.logs.Log("reconnected to photo service successfully")
	return nil
}

func (a *photoAdapter) CalculatePhotoPrice(ctx context.Context, userId, creatorId string, photoIds []string) (*[]*model.CheckoutItem, *model.Total, error) {
	processPhotoRequest := &photopb.CalculatePhotoPriceRequest{
		UserId:    userId,
		CreatorId: creatorId,
		PhotoIds:  photoIds,
	}

	const maxAttempts = 3

	var response *photopb.CalculatePhotoPriceResponse

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		res, err := a.client.CalculatePhotoPrice(ctx, processPhotoRequest)
		if err == nil {
			response = res
			break
		}

		a.logs.Error(fmt.Sprintf("attempt %d: failed to calculate photo price: %v", attempt, err))
		if attempt < maxAttempts {
			if recErr := a.reconnect(ctx); recErr != nil {
				a.logs.Error(fmt.Sprintf("attempt %d: reconnect failed: %v", attempt, recErr))
			}
		} else {
			return nil, nil, helper.FromGRPCError(err)
		}
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

	const maxAttempts = 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		_, err := a.client.OwnerOwnPhotos(ctx, ownerOwnPhotosRequest)
		if err == nil {
			break
		}

		a.logs.Error(fmt.Sprintf("attempt %d: failed to calculate photo price: %v", attempt, err))
		if attempt < maxAttempts {
			if recErr := a.reconnect(ctx); recErr != nil {
				a.logs.Error(fmt.Sprintf("attempt %d: reconnect failed: %v", attempt, recErr))
			}
		} else {
			return helper.FromGRPCError(err)
		}
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

func (a *photoAdapter) GetCreator(ctx context.Context, userId string) (*entity.Creator, error) {
	pbRequest := &photopb.GetCreatorRequest{
		UserId: userId,
	}

	pbResponse, err := a.client.GetCreator(context.Background(), pbRequest)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	creator := &entity.Creator{
		Id: pbResponse.GetCreator().GetId(),
	}

	return creator, nil
}

func (a *photoAdapter) CalculatePhotoPriceV2(ctx context.Context, userId, creatorId string, request *model.CreateTransactionV2Request) (*[]*model.CheckoutItem, *model.Total, error) {
	checkoutItemPb := make([]*photopb.CheckoutItemWeb, 0, len(request.Items))
	for _, item := range request.Items {
		var discount *photopb.Discount
		if item.Discount != nil {
			discount = &photopb.Discount{
				Discount:            item.Discount.Discount,
				DiscountId:          item.Discount.DiscountId,
				DiscountType:        item.Discount.DiscountType,
				DiscountMinQuantity: int32(item.Discount.DiscountMinQuantity),
				DiscountValue:       item.Discount.DiscountValue,
			}
		}

		checkoutItemPb = append(checkoutItemPb,
			&photopb.CheckoutItemWeb{
				PhotoId:    item.PhotoId,
				CreatorId:  item.CreatorId,
				Title:      item.Title,
				Price:      item.Price,
				Discount:   discount,
				FinalPrice: item.FinalPrice,
			})
	}

	processPhotoRequest := &photopb.CalculatePhotoPriceV2Request{
		UserId:         userId,
		CreatorId:      creatorId,
		ChekoutItemWeb: checkoutItemPb,
		TotalPrice:     request.TotalPrice,
		TotalDiscount:  request.TotalDiscount,
	}

	response, err := a.client.CalculatePhotoPriceV2(ctx, processPhotoRequest)
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
