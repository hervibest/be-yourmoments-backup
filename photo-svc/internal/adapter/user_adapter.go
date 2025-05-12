package adapter

import (
	"context"
	"log"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/discovery"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/utils"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
	userpb "github.com/hervibest/be-yourmoments-backup/pb/user"
)

type UserAdapter interface {
	AuthenticateUser(ctx context.Context, token string) (*userpb.AuthenticateResponse, error)
	SendBulkPhotoNotification(ctx context.Context, request []*photopb.BulkUserSimilarPhoto) (*userpb.SendBulkPhotoNotificationResponse, error)
	SendSinglePhotoNotification(ctx context.Context, request []*photopb.UserSimilarPhoto) (*userpb.SendSinglePhotoNotificationResponse, error)
	SendBulkNotification(ctx context.Context, countMap map[string]int32) (*userpb.SendBulkNotificationResponse, error)
	SendSingleFacecamNotificaton(ctx context.Context, request []*photopb.UserSimilarPhoto) (*userpb.SendSingleFacecamNotificationResponse, error)
}

type userAdapter struct {
	client userpb.UserServiceClient
}

func NewUserAdapter(ctx context.Context, registry discovery.Registry) (UserAdapter, error) {
	userServiceName := utils.GetEnv("USER_SVC_NAME")
	conn, err := discovery.ServiceConnection(ctx, userServiceName, registry)
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

func (a *userAdapter) SendSinglePhotoNotification(ctx context.Context, request []*photopb.UserSimilarPhoto) (*userpb.SendSinglePhotoNotificationResponse, error) {
	sendSinglePhotoNotificationReq := &userpb.SendSinglePhotoNotificationRequest{
		UserSimilarPhoto: request,
	}

	response, err := a.client.SendSinglePhotoNotification(ctx, sendSinglePhotoNotificationReq)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	return response, nil
}

func (a *userAdapter) SendSingleFacecamNotificaton(ctx context.Context, request []*photopb.UserSimilarPhoto) (*userpb.SendSingleFacecamNotificationResponse, error) {
	sendSingleFacecamNotificationReq := &userpb.SendSingleFacecamNotificationRequest{
		UserSimilarPhoto: request,
	}

	response, err := a.client.SendSingleFacecamNotification(ctx, sendSingleFacecamNotificationReq)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	return response, nil
}

// THIS IS THE NEWEST VERSION USING COUNT IN PHOTO SVC
func (a *userAdapter) SendBulkNotification(ctx context.Context, countMap map[string]int32) (*userpb.SendBulkNotificationResponse, error) {
	sendBulkPhotoNotificationReq := &userpb.SendBulkNotificationRequest{
		CountMap: countMap,
	}

	response, err := a.client.SendBulkNotification(ctx, sendBulkPhotoNotificationReq)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	return response, nil
}
