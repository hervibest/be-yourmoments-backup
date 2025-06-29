package grpc

import (
	"context"
	"log"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"google.golang.org/grpc/codes"
)

func (h *PhotoGRPCHandler) CalculatePhotoPrice(ctx context.Context, pbReq *photopb.CalculatePhotoPriceRequest) (
	*photopb.CalculatePhotoPriceResponse, error) {
	log.Println("----  Calcualte Photo Price Requets via GRPC in photo-svc ------")

	request := &model.CalculateRequest{
		UserId:    pbReq.GetUserId(),
		CreatorId: pbReq.GetCreatorId(),
		PhotoIds:  pbReq.GetPhotoIds(),
	}

	items, total, err := h.checkoutUseCase.LockPhotosAndCalculatePrice(context.Background(), request)
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

	request := &model.OwnerOwnPhotosRequest{
		OwnerId:  pbReq.GetOwnerId(),
		PhotoIds: pbReq.GetPhotoIds(),
	}

	if err := h.checkoutUseCase.OwnerOwnPhotos(context.Background(), request); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &photopb.OwnerOwnPhotosResponse{
		Status: int64(codes.OK),
	}, nil
}

func (h *PhotoGRPCHandler) CancelPhotos(ctx context.Context, pbReq *photopb.CancelPhotosRequest) (
	*photopb.CancelPhotosResponse, error) {
	log.Println("----  Cancel Photos Request via GRPC in photo-svc ------")

	request := &model.CancelPhotosRequest{
		UserId:   pbReq.GetUserId(),
		PhotoIds: pbReq.GetPhotoIds(),
	}
	if err := h.checkoutUseCase.CancelPhotos(context.Background(), request); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &photopb.CancelPhotosResponse{
		Status: int64(codes.OK),
	}, nil
}

func (h *PhotoGRPCHandler) CalculatePhotoPriceV2(ctx context.Context, pbReq *photopb.CalculatePhotoPriceV2Request) (
	*photopb.CalculatePhotoPriceV2Response, error) {
	log.Println("----  Calcualte Photo Price Requets via GRPC in photo-svc ------")

	checkoutItemWeb := make([]model.CheckoutItemWeb, 0, len(pbReq.GetChekoutItemWeb()))
	for _, item := range pbReq.GetChekoutItemWeb() {
		var discount *model.DiscountItem
		if item.GetDiscount() != nil {
			discount = &model.DiscountItem{
				Discount:            item.GetDiscount().GetDiscount(),
				DiscountMinQuantity: int(item.GetDiscount().GetDiscountMinQuantity()),
				DiscountValue:       item.GetDiscount().GetDiscountValue(),
				DiscountId:          item.GetDiscount().GetDiscountId(),
				DiscountType:        enum.DiscountType(item.GetDiscount().GetDiscountType()),
			}
		}
		checkoutItemWeb = append(checkoutItemWeb,
			model.CheckoutItemWeb{
				PhotoId:    item.GetPhotoId(),
				CreatorId:  item.GetCreatorId(),
				Title:      item.GetTitle(),
				Discount:   discount,
				Price:      item.GetPrice(),
				FinalPrice: item.GetFinalPrice(),
			})
	}

	request := &model.CalculateV2Request{
		UserId:        pbReq.GetUserId(),
		CreatorId:     pbReq.GetCreatorId(),
		Items:         checkoutItemWeb,
		TotalPrice:    pbReq.GetTotalPrice(),
		TotalDiscount: pbReq.GetTotalDiscount(),
	}

	items, total, err := h.checkoutUseCase.LockPhotosAndCalculatePriceV2(context.Background(), request)
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

	return &photopb.CalculatePhotoPriceV2Response{
		Status: int64(codes.OK),
		Items:  pbItemReponses,
		Total:  totalPbResponse,
	}, nil
}
