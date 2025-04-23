package grpc

import (
	"be-yourmoments/user-svc/internal/helper"
	"be-yourmoments/user-svc/internal/model"
	"be-yourmoments/user-svc/internal/usecase"
	"context"
	"log"

	"github.com/be-yourmoments/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type PhotoGRPCHandler struct {
	usecase usecase.AuthUseCase
	pb.UnimplementedUserServiceServer
}

func NewPhotoGRPCHandler(server *grpc.Server, usecase usecase.AuthUseCase) {
	handler := &PhotoGRPCHandler{
		usecase: usecase,
	}

	pb.RegisterUserServiceServer(server, handler)
}

func (h *PhotoGRPCHandler) Authenticate(ctx context.Context,
	pbReq *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {

	log.Println("---- Authenticate User via gRPC in user-svc ------")

	request := &model.VerifyUserRequest{
		Token: pbReq.GetToken(),
	}

	response, err := h.usecase.Verify(ctx, request)
	if err != nil {
		return nil, helper.ErrGRPC(err)
	}

	userPb := &pb.User{
		UserId:      response.UserId,
		Username:    response.Username,
		Email:       response.Email,
		PhoneNumber: response.PhoneNumber,
		CreatorId:   response.CreatorId,
		WalletId:    response.WalletId,
	}

	return &pb.AuthenticateResponse{
		Status: int64(codes.OK),
		User:   userPb,
	}, nil
}
