package grpc

import (
	"context"
	"log"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
