package grpc

import (
	"context"
	"log"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"google.golang.org/grpc/codes"
)

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
