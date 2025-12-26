package adapter

import (
	"context"
	"log"
	"time"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/discovery"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/utils"
	"google.golang.org/protobuf/types/known/timestamppb"

	userpb "github.com/hervibest/be-yourmoments-backup/pb/user"
)

type UserAdapter interface {
	AuthenticateUser(ctx context.Context, token string) (*userpb.AuthenticateResponse, error)
	AuthenticateUserV2(ctx context.Context, token, userId string, expiresAt time.Time) (*userpb.AuthenticateResponseV2, error)
}

type userAdapter struct {
	client userpb.UserServiceClient
}

func NewUserAdapter(ctx context.Context, registry discovery.Registry, logs *logger.Log) (UserAdapter, error) {
	userServiceName := utils.GetEnv("USER_SVC_NAME")
	conn, err := discovery.ServiceConnection(ctx, userServiceName, registry, logs)
	if err != nil {
		return nil, err
	}

	log.Print("successfuly connected to user-svc-grpc")
	client := userpb.NewUserServiceClient(conn)

	return &userAdapter{
		client: client,
	}, nil
}

func (a *userAdapter) AuthenticateUser(ctx context.Context, token string) (*userpb.AuthenticateResponse, error) {
	processPhotoRequest := &userpb.AuthenticateRequest{
		Token: token,
	}

	response, err := a.client.Authenticate(ctx, processPhotoRequest)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	return response, nil
}

func (a *userAdapter) AuthenticateUserV2(ctx context.Context, token, userId string, expiresAt time.Time) (*userpb.AuthenticateResponseV2, error) {
	processPhotoRequest := &userpb.AuthenticateRequestV2{
		Token:     token,
		UserId:    userId,
		ExpiresAt: timestamppb.New(expiresAt),
	}

	response, err := a.client.AuthenticateV2(ctx, processPhotoRequest)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	return response, nil
}
