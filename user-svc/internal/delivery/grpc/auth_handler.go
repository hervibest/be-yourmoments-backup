package grpc

import (
	"be-yourmoments/user-svc/internal/helper"
	"be-yourmoments/user-svc/internal/model"
	"context"
	"log"

	userpb "github.com/be-yourmoments/pb/user"

	"google.golang.org/grpc/codes"
)

func (h *UserGRPCHandler) Authenticate(ctx context.Context,
	pbReq *userpb.AuthenticateRequest) (*userpb.AuthenticateResponse, error) {

	log.Println("---- Authenticate User via gRPC in user-svc ------")

	request := &model.VerifyUserRequest{
		Token: pbReq.GetToken(),
	}

	response, err := h.usecase.Verify(ctx, request)
	if err != nil {
		return nil, helper.ErrGRPC(err)
	}

	userPb := &userpb.User{
		UserId:      response.UserId,
		Username:    response.Username,
		Email:       response.Email,
		PhoneNumber: response.PhoneNumber,
		CreatorId:   response.CreatorId,
		WalletId:    response.WalletId,
	}

	return &userpb.AuthenticateResponse{
		Status: int64(codes.OK),
		User:   userPb,
	}, nil
}
