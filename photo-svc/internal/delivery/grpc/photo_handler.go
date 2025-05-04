package grpc

import (
	"context"
	"log"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"

	"google.golang.org/grpc/codes"
)

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

func (h *PhotoGRPCHandler) GetPhotoWithDetails(ctx context.Context, pbReq *photopb.GetPhotoWithDetailsRequest) (
	*photopb.GetPhotoWithDetailsResponse, error) {
	log.Println("----  CreatePhoto Bulk Photo Requets via GRPC in photo-svc ------")
	response, err := h.photoUseCase.UserGetPhotoWithDetail(context.Background(), pbReq.GetPhotoIds(), pbReq.GetUserId())
	if err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return response, nil
}
