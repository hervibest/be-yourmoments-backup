package adapter

import (
	"context"
	"log"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/discovery"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
	userpb "github.com/hervibest/be-yourmoments-backup/pb/user"
)

type UserAdapter interface {
	AuthenticateUser(ctx context.Context, token string) (*userpb.AuthenticateResponse, error)
	SendBulkPhotoNotification(ctx context.Context, request []*photopb.BulkUserSimilarPhoto) (*userpb.SendBulkPhotoNotificationResponse, error)
}

type userAdapter struct {
	client userpb.UserServiceClient
}

func NewUserAdapter(ctx context.Context, registry discovery.Registry) (UserAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "user-svc-grpc", registry)
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

func (a *userAdapter) SendBulkPhotoNotification(ctx context.Context, request []*photopb.BulkUserSimilarPhoto) (*userpb.SendBulkPhotoNotificationResponse, error) {
	sendBulkPhotoNotificationReq := &userpb.SendBulkPhotoNotificationRequest{
		BulkUserSimilarPhoto: request,
	}

	response, err := a.client.SendBulkPhotoNotification(ctx, sendBulkPhotoNotificationReq)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	return response, nil
}
