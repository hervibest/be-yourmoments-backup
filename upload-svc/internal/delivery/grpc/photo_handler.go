package grpc

import (
	"context"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PhotoGRPCHandler struct {
	usecase usecase.PhotoUsecase
	photopb.UnimplementedPhotoServiceServer
}

func NewPhotoGRPCHandler(server *grpc.Server, usecase usecase.PhotoUsecase) {
	handler := &PhotoGRPCHandler{
		usecase: usecase,
	}

	photopb.RegisterPhotoServiceServer(server, handler)
}

func (h *PhotoGRPCHandler) UpdatePhotographerPhoto(ctx context.Context,
	pbReq *photopb.UpdatePhotographerPhotoRequest) (
	*photopb.UpdatePhotographerPhotoResponse, error) {

	// req := converter.GrpcToCreateRequest(pbReq)
	// h.usecase.UpdatePhoto(ctx, req)

	return nil, nil
}
func (h *PhotoGRPCHandler) UpdateFaceRecogPhoto(ctx context.Context, req *photopb.UpdateFaceRecogPhotoRequest) (*photopb.UpdateFaceRecogPhotoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFaceRecogPhoto not implemented")
}
