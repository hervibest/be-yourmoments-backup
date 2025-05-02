package grpc

import (
	"context"
	"log"

	userpb "github.com/hervibest/be-yourmoments-backup/pb/user"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper"

	"google.golang.org/grpc/codes"
)

func (h *UserGRPCHandler) SendBulkPhotoNotification(ctx context.Context,
	pbReq *userpb.SendBulkPhotoNotificationRequest) (*userpb.SendBulkPhotoNotificationResponse, error) {

	log.Println("---- Send Bulk Photo Notification User via gRPC in user-svc ------")

	if err := h.notificationUseCase.ProcessAndSendBulkNotifications(ctx, pbReq.GetBulkUserSimilarPhoto()); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &userpb.SendBulkPhotoNotificationResponse{
		Status: int64(codes.OK),
	}, nil
}

func (h *UserGRPCHandler) SendBulkNotification(ctx context.Context,
	pbReq *userpb.SendBulkNotificationRequest) (*userpb.SendBulkNotificationResponse, error) {

	log.Println("---- Send Bulk Photo Notification (Newest Version) User via gRPC in user-svc ------")

	if err := h.notificationUseCase.ProcessAndSendBulkNotificationsV2(ctx, pbReq); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &userpb.SendBulkNotificationResponse{
		Status: int64(codes.OK),
	}, nil
}

func (h *UserGRPCHandler) SendSinglePhotoNotification(ctx context.Context, pbReq *userpb.SendSinglePhotoNotificationRequest) (*userpb.SendSinglePhotoNotificationResponse, error) {

	log.Println("---- Send Single Photo Notification User via gRPC in user-svc ------")

	if err := h.notificationUseCase.ProcessAndSendSingleNotifications(ctx, pbReq.GetUserSimilarPhoto()); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &userpb.SendSinglePhotoNotificationResponse{
		Status: int64(codes.OK),
	}, nil
}

// TODO (URGENT) FIX API CONTRACT SendBulk is in PhotoContract not user contract
func (h *UserGRPCHandler) SendSingleFacecamNotification(ctx context.Context, pbReq *userpb.SendSingleFacecamNotificationRequest) (*userpb.SendSingleFacecamNotificationResponse, error) {

	log.Println("---- Send Single Facecam Notification User via gRPC in user-svc ------")

	if err := h.notificationUseCase.ProcessAndSendSingleFacecamNotifications(ctx, pbReq.GetUserSimilarPhoto()); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &userpb.SendSingleFacecamNotificationResponse{
		Status: int64(codes.OK),
	}, nil
}
