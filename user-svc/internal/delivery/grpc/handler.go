package grpc

import (
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/usecase"

	userpb "github.com/hervibest/be-yourmoments-backup/pb/user"
	"google.golang.org/grpc"
)

type UserGRPCHandler struct {
	authUseCase         usecase.AuthUseCase
	notificationUseCase usecase.NotificationUseCase
	userpb.UnimplementedUserServiceServer
}

func NewUserGRPCHandler(server *grpc.Server, authUseCase usecase.AuthUseCase, notificationUseCase usecase.NotificationUseCase) {
	handler := &UserGRPCHandler{
		authUseCase:         authUseCase,
		notificationUseCase: notificationUseCase,
	}

	userpb.RegisterUserServiceServer(server, handler)
}
