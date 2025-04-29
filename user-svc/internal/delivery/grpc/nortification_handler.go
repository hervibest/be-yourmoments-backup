package grpc

import (
	"context"
	"log"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper"

	"google.golang.org/grpc/codes"
)

// TODO (URGENT) FIX API CONTRACT SendBulk is in PhotoContract not user contract
func (h *UserGRPCHandler) SendBulkPhotoNotification(ctx context.Context,
	pbReq *photopb.SendBulkPhotoNotificationRequest) (*pb.SendBulkPhotoNotificationResponse, error) {

	log.Println("---- Send Bulk Photo Notification User via gRPC in user-svc ------")

	if err := h.usecase.ProcessAndSendNotifications(ctx, pbReq.GetBulkUserSimilarPhoto()); err != nil {
		return nil, helper.ErrGRPC(err)
	}

	return &photopb.SendBulkPhotoNotificationResponse{
		Status: int64(codes.OK),
	}, nil
}
