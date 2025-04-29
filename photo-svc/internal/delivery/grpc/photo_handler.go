package grpc

import (
	"context"
	"log"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/usecase"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PhotoGRPCHandler struct {
	photoUseCase            usecase.PhotoUseCase
	facecamUseCase          usecase.FacecamUseCase
	userSimilarPhotoUseCase usecase.UserSimilarUsecase
	creatorUseCase          usecase.CreatorUseCase
	checkoutUseCase         usecase.CheckoutUseCase

	photopb.UnimplementedPhotoServiceServer
}

func NewPhotoGRPCHandler(server *grpc.Server, photoUseCase usecase.PhotoUseCase,
	facecamUseCase usecase.FacecamUseCase, userSimilarPhotoUseCase usecase.UserSimilarUsecase,
	creatorUseCase usecase.CreatorUseCase, checkoutUseCase usecase.CheckoutUseCase) {
	handler := &PhotoGRPCHandler{
		photoUseCase:            photoUseCase,
		facecamUseCase:          facecamUseCase,
		userSimilarPhotoUseCase: userSimilarPhotoUseCase,
		creatorUseCase:          creatorUseCase,
		checkoutUseCase:         checkoutUseCase,
	}

	photopb.RegisterPhotoServiceServer(server, handler)
}

func (h *PhotoGRPCHandler) CreatePhoto(ctx context.Context, pbReq *photopb.CreatePhotoRequest) (
	*photopb.CreatePhotoResponse, error) {
	log.Println("----  CreatePhoto Requets via GRPC in photo-svc ------")
	if err := h.photoUseCase.CreatePhoto(context.Background(), pbReq); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &photopb.CreatePhotoResponse{
		Status: int64(codes.OK),
	}, nil
}

func (h *PhotoGRPCHandler) CreateUserSimilar(ctx context.Context, pbReq *photopb.CreateUserSimilarPhotoRequest) (
	*photopb.CreateUserSimilarPhotoResponse, error) {
	log.Println("----  CreatePhoto user similar Requets via GRPC in photo-svc ------")
	if err := h.userSimilarPhotoUseCase.CreateUserSimilar(context.Background(), pbReq); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &photopb.CreateUserSimilarPhotoResponse{
		Status: int64(codes.OK),
	}, nil
}

func (h *PhotoGRPCHandler) UpdatePhotoDetail(ctx context.Context, pbReq *photopb.UpdatePhotoDetailRequest) (
	*photopb.UpdatePhotoDetailResponse, error) {
	log.Println("----  UpdatePhoto Requets via GRPC in photo-svc ------")
	if err := h.photoUseCase.UpdatePhotoDetail(context.Background(), pbReq); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &photopb.UpdatePhotoDetailResponse{
		Status: int64(codes.OK),
	}, nil
}

func (h *PhotoGRPCHandler) CreateFacecam(ctx context.Context, pbReq *photopb.CreateFacecamRequest) (
	*photopb.CreateFacecamResponse, error) {
	log.Println("----  Create facecam Requets via GRPC in photo-svc ------")
	if err := h.facecamUseCase.CreateFacecam(context.Background(), pbReq); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &photopb.CreateFacecamResponse{
		Status: int64(codes.OK),
	}, nil
}

func (h *PhotoGRPCHandler) CreateUserSimilarFacecam(ctx context.Context, pbReq *photopb.CreateUserSimilarFacecamRequest) (
	*photopb.CreateUserSimilarFacecamResponse, error) {
	log.Println("----  CreatePhoto user similar Requets via GRPC in photo-svc ------")
	if err := h.userSimilarPhotoUseCase.CreateUserFacecam(context.Background(), pbReq); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &photopb.CreateUserSimilarFacecamResponse{
		Status: int64(codes.OK),
	}, nil
}

func (h *PhotoGRPCHandler) CreateCreator(ctx context.Context, pbReq *photopb.CreateCreatorRequest) (
	*photopb.CreateCreatorResponse, error) {
	log.Println("----  CreatePhoto user similar Requets via GRPC in photo-svc ------")

	request := &model.CreateCreatorRequest{
		UserId: pbReq.GetUserId(),
	}

	response, err := h.creatorUseCase.CreateCreator(context.Background(), request)
	if err != nil {
		return nil, helper.ErrGRPC(err)
	}

	creatorPb := &photopb.Creator{
		Id:     response.Id,
		UserId: response.UserId,
		CreatedAt: &timestamppb.Timestamp{
			Seconds: int64(response.CreatedAt.Second()),
			Nanos:   int32(response.CreatedAt.UnixNano()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: int64(response.UpdatedAt.Second()),
			Nanos:   int32(response.UpdatedAt.UnixNano()),
		},
	}

	return &photopb.CreateCreatorResponse{
		Status:  int64(codes.OK),
		Creator: creatorPb,
	}, nil
}

func (h *PhotoGRPCHandler) GetCreator(ctx context.Context, pbReq *photopb.GetCreatorRequest) (
	*photopb.GetCreatorResponse, error) {
	log.Println("----  GetCreator Requets via GRPC in photo-svc ------")

	request := &model.GetCreatorRequest{
		UserId: pbReq.GetUserId(),
	}

	response, err := h.creatorUseCase.GetCreator(context.Background(), request)
	if err != nil {
		return nil, helper.ErrGRPC(err)
	}

	creatorPb := &photopb.Creator{
		Id:     response.Id,
		UserId: response.UserId,
		CreatedAt: &timestamppb.Timestamp{
			Seconds: int64(response.CreatedAt.Second()),
			Nanos:   int32(response.CreatedAt.UnixNano()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: int64(response.UpdatedAt.Second()),
			Nanos:   int32(response.UpdatedAt.UnixNano()),
		},
	}

	return &photopb.GetCreatorResponse{
		Status:  int64(codes.OK),
		Creator: creatorPb,
	}, nil
}

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

func (h *PhotoGRPCHandler) CreateBulkPhoto(ctx context.Context, pbReq *photopb.CreateBulkPhotoRequest) (
	*photopb.CreateBulkPhotoResponse, error) {
	log.Println("----  CreatePhoto Bulk Photo Requets via GRPC in photo-svc ------")
	if err := h.photoUseCase.CreateBulkPhoto(context.Background(), pbReq); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &photopb.CreateBulkPhotoResponse{
		Status: int64(codes.OK),
	}, nil
}

func (h *PhotoGRPCHandler) CreateBulkUserSimilarPhotos(ctx context.Context, pbReq *photopb.CreateBulkUserSimilarPhotoRequest) (
	*photopb.CreateBulkUserSimilarPhotoResponse, error) {
	log.Println("----  Create Bulk User Similar Photos Requets via GRPC in photo-svc ------")
	if err := h.userSimilarPhotoUseCase.CreateBulkUserSimilarPhotos(context.Background(), pbReq); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &photopb.CreateBulkUserSimilarPhotoResponse{
		Status: int64(codes.OK),
	}, nil
}
