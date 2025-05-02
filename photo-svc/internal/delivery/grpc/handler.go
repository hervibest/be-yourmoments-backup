package grpc

import (
	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/usecase"
	"google.golang.org/grpc"
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
