package grpc

import (
	"be-yourmoments/photo-svc/internal/helper"
	"be-yourmoments/photo-svc/internal/model"
	"be-yourmoments/photo-svc/internal/usecase"
	"context"
	"log"
	"net/http"

	"github.com/be-yourmoments/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PhotoGRPCHandler struct {
	photoUseCase            usecase.PhotoUsecase
	facecamUseCase          usecase.FacecamUseCase
	userSimilarPhotoUseCase usecase.UserSimilarUsecase
	creatorUseCase          usecase.CreatorUseCase
	checkoutUseCase         usecase.CheckoutUseCase

	pb.UnimplementedPhotoServiceServer
}

func NewPhotoGRPCHandler(server *grpc.Server, photoUseCase usecase.PhotoUsecase,
	facecamUseCase usecase.FacecamUseCase, userSimilarPhotoUseCase usecase.UserSimilarUsecase,
	creatorUseCase usecase.CreatorUseCase, checkoutUseCase usecase.CheckoutUseCase) {
	handler := &PhotoGRPCHandler{
		photoUseCase:            photoUseCase,
		facecamUseCase:          facecamUseCase,
		userSimilarPhotoUseCase: userSimilarPhotoUseCase,
		creatorUseCase:          creatorUseCase,
		checkoutUseCase:         checkoutUseCase,
	}

	pb.RegisterPhotoServiceServer(server, handler)
}

func (h *PhotoGRPCHandler) CreatePhoto(ctx context.Context, pbReq *pb.CreatePhotoRequest) (
	*pb.CreatePhotoResponse, error) {
	log.Println("----  CreatePhoto Requets via GRPC in photo-svc ------")
	if err := h.photoUseCase.CreatePhoto(context.Background(), pbReq); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &pb.CreatePhotoResponse{
		Status: http.StatusCreated,
	}, nil
}

func (h *PhotoGRPCHandler) CreateUserSimilar(ctx context.Context, pbReq *pb.CreateUserSimilarPhotoRequest) (
	*pb.CreateUserSimilarPhotoResponse, error) {
	log.Println("----  CreatePhoto user similar Requets via GRPC in photo-svc ------")
	if err := h.userSimilarPhotoUseCase.CreateUserSimilar(context.Background(), pbReq); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &pb.CreateUserSimilarPhotoResponse{
		Status: int64(codes.OK),
	}, nil
}

func (h *PhotoGRPCHandler) UpdatePhotoDetail(ctx context.Context, pbReq *pb.UpdatePhotoDetailRequest) (
	*pb.UpdatePhotoDetailResponse, error) {
	log.Println("----  UpdatePhoto Requets via GRPC in photo-svc ------")
	if err := h.photoUseCase.UpdatePhotoDetail(context.Background(), pbReq); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &pb.UpdatePhotoDetailResponse{
		Status: int64(codes.OK),
	}, nil
}

func (h *PhotoGRPCHandler) CreateFacecam(ctx context.Context, pbReq *pb.CreateFacecamRequest) (
	*pb.CreateFacecamResponse, error) {
	log.Println("----  Create facecam Requets via GRPC in photo-svc ------")
	if err := h.facecamUseCase.CreateFacecam(context.Background(), pbReq); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &pb.CreateFacecamResponse{
		Status: int64(codes.OK),
	}, nil
}

func (h *PhotoGRPCHandler) CreateUserSimilarFacecam(ctx context.Context, pbReq *pb.CreateUserSimilarFacecamRequest) (
	*pb.CreateUserSimilarFacecamResponse, error) {
	log.Println("----  CreatePhoto user similar Requets via GRPC in photo-svc ------")
	if err := h.userSimilarPhotoUseCase.CreateUserFacecam(context.Background(), pbReq); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &pb.CreateUserSimilarFacecamResponse{
		Status: int64(codes.OK),
	}, nil
}

func (h *PhotoGRPCHandler) CreateCreator(ctx context.Context, pbReq *pb.CreateCreatorRequest) (
	*pb.CreateCreatorResponse, error) {
	log.Println("----  CreatePhoto user similar Requets via GRPC in photo-svc ------")

	request := &model.CreateCreatorRequest{
		UserId: pbReq.GetUserId(),
	}

	response, err := h.creatorUseCase.CreateCreator(context.Background(), request)
	if err != nil {
		return nil, helper.ErrGRPC(err)
	}

	creatorPb := &pb.Creator{
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

	return &pb.CreateCreatorResponse{
		Status:  int64(codes.OK),
		Creator: creatorPb,
	}, nil
}

func (h *PhotoGRPCHandler) CalculatePhotoPrice(ctx context.Context, pbReq *pb.CalculatePhotoPriceRequest) (
	*pb.CalculatePhotoPriceResponse, error) {
	log.Println("----  Calcualte Photo Price Requets via GRPC in photo-svc ------")

	request := &model.PreviewCheckoutRequest{
		UserId:   pbReq.GetUserId(),
		PhotoIds: pbReq.GetPhotoIds(),
	}

	items, total, err := h.checkoutUseCase.CalculatePrice(context.Background(), request)
	if err != nil {
		return nil, helper.ErrGRPC(err)
	}

	pbItemReponses := make([]*pb.CheckoutItem, 0)
	for _, item := range *items {
		pbResponse := &pb.CheckoutItem{
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

	totalPbResponse := &pb.Total{
		Price:    total.Price,
		Discount: total.Discount,
	}

	return &pb.CalculatePhotoPriceResponse{
		Status: int64(codes.OK),
		Items:  pbItemReponses,
		Total:  totalPbResponse,
	}, nil
}

func (h *PhotoGRPCHandler) OwnerOwnPhotos(ctx context.Context, pbReq *pb.OwnerOwnPhotosRequest) (
	*pb.OwnerOwnPhotosResponse, error) {
	log.Println("----  OwnerOwnPhotos Requets via GRPC in photo-svc ------")
	if err := h.checkoutUseCase.OwnerOwnPhotos(context.Background(), pbReq.GetOwnerId(), pbReq.GetPhotoIds()); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &pb.OwnerOwnPhotosResponse{
		Status: int64(codes.OK),
	}, nil
}
