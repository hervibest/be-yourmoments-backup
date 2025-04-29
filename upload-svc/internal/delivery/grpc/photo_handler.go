package grpc

import (
	"be-yourmoments/upload-svc/internal/usecase"
	"context"

	photopb "github.com/be-yourmoments/pb/photo"

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
