package grpc

import (
	"context"
	"log"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model"

	userpb "github.com/hervibest/be-yourmoments-backup/pb/user"

	"google.golang.org/grpc/codes"
)

func (h *UserGRPCHandler) Authenticate(ctx context.Context,
	pbReq *userpb.AuthenticateRequest) (*userpb.AuthenticateResponse, error) {

	log.Println("---- Authenticate User via gRPC in user-svc ------")

	request := &model.VerifyUserRequest{
		Token: pbReq.GetToken(),
	}

	response, err := h.authUseCase.Verify(ctx, request)
	if err != nil {
		return nil, helper.ErrGRPC(err)
	}

	userPb := &userpb.User{
		UserId:      response.UserId,
		Username:    response.Username,
		Email:       response.Email,
		PhoneNumber: response.PhoneNumber,
		Similarity:  uint32(response.Similarity),
		CreatorId:   response.CreatorId,
		WalletId:    response.WalletId,
	}

	return &userpb.AuthenticateResponse{
		Status: int64(codes.OK),
		User:   userPb,
	}, nil
}
