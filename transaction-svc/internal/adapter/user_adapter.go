package adapter

import (
	"be-yourmoments/transaction-svc/internal/helper"
	"be-yourmoments/transaction-svc/internal/helper/discovery"
	"context"
	"log"

	"github.com/be-yourmoments/pb"
)

type UserAdapter interface {
	AuthenticateUser(ctx context.Context, token string) (*pb.AuthenticateResponse, error)
}

type userAdapter struct {
	client pb.UserServiceClient
}

func NewUserAdapter(ctx context.Context, registry discovery.Registry) (UserAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "user-svc-grpc", registry)
	if err != nil {
		return nil, err
	}

	log.Print("successfuly connected to user-svc-grpc")
	client := pb.NewUserServiceClient(conn)

	return &userAdapter{
		client: client,
	}, nil
}

func (a *userAdapter) AuthenticateUser(ctx context.Context, token string) (*pb.AuthenticateResponse, error) {
	processPhotoRequest := &pb.AuthenticateRequest{
		Token: token,
	}

	response, err := a.client.Authenticate(ctx, processPhotoRequest)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	return response, nil
}
