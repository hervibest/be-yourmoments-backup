package grpc

import (
	"context"
	"log"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"google.golang.org/grpc/codes"
)

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
