package grpc

import (
	"context"
	"log"

	userpb "github.com/hervibest/be-yourmoments-backup/pb/user"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper"

	"google.golang.org/grpc/codes"
)

// TODO (URGENT) FIX API CONTRACT SendBulk is in PhotoContract not user contract
func (h *UserGRPCHandler) SendBulkPhotoNotification(ctx context.Context,
	pbReq *userpb.SendBulkPhotoNotificationRequest) (*userpb.SendBulkPhotoNotificationResponse, error) {

	log.Println("---- Send Bulk Photo Notification User via gRPC in user-svc ------")

	if err := h.notificationUseCase.ProcessAndSendNotifications(ctx, pbReq.GetBulkUserSimilarPhoto()); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &userpb.SendBulkPhotoNotificationResponse{
		Status: int64(codes.OK),
	}, nil
}
