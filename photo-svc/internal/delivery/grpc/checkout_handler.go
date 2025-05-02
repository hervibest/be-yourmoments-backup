package grpc

import (
	"context"
	"log"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"google.golang.org/grpc/codes"
)

func (h *PhotoGRPCHandler) CalculatePhotoPrice(ctx context.Context, pbReq *photopb.CalculatePhotoPriceRequest) (
	*photopb.CalculatePhotoPriceResponse, error) {
	log.Println("----  Calcualte Photo Price Requets via GRPC in photo-svc ------")

	request := &model.PreviewCheckoutRequest{
		UserId:   pbReq.GetUserId(),
		PhotoIds: pbReq.GetPhotoIds(),
	}

	items, total, err := h.checkoutUseCase.CalculatePrice(context.Background(), request)
	if err != nil {
		return nil, helper.ErrGRPC(err)
	}

	pbItemReponses := make([]*photopb.CheckoutItem, 0)
	for _, item := range *items {
		pbResponse := &photopb.CheckoutItem{
			PhotoId:             item.PhotoId,
			CreatorId:           item.CreatorId,
			Title:               item.Title,
			YourMomentsUrl:      item.YourMomentsUrl,
			Price:               item.Price,
			Discount:            item.Discount,
			DiscountValue:       item.DiscountValue,
			DiscountMinQuantity: int32(item.DiscountMinQuantity),
			DiscountId:          item.DiscountId,
			DiscountType:        string(item.DiscountType),
			FinalPrice:          item.FinalPrice,
		}

		pbItemReponses = append(pbItemReponses, pbResponse)
	}

	totalPbResponse := &photopb.Total{
		Price:    total.Price,
		Discount: total.Discount,
	}

	return &photopb.CalculatePhotoPriceResponse{
		Status: int64(codes.OK),
		Items:  pbItemReponses,
		Total:  totalPbResponse,
	}, nil
}

func (h *PhotoGRPCHandler) OwnerOwnPhotos(ctx context.Context, pbReq *photopb.OwnerOwnPhotosRequest) (
	*photopb.OwnerOwnPhotosResponse, error) {
	log.Println("----  OwnerOwnPhotos Requets via GRPC in photo-svc ------")
	if err := h.checkoutUseCase.OwnerOwnPhotos(context.Background(), pbReq.GetOwnerId(), pbReq.GetPhotoIds()); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &photopb.OwnerOwnPhotosResponse{
		Status: int64(codes.OK),
	}, nil
}
