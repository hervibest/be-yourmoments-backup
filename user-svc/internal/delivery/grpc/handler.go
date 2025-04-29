package grpc

import (
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/usecase"

	userpb "github.com/hervibest/be-yourmoments-backup/pb/user"
	"google.golang.org/grpc"
)

type UserGRPCHandler struct {
	usecase usecase.AuthUseCase
	userpb.UnimplementedUserServiceServer
}

func NewUserGRPCHandler(server *grpc.Server, usecase usecase.AuthUseCase) {
	handler := &UserGRPCHandler{
		usecase: usecase,
	}

	userpb.RegisterUserServiceServer(server, handler)
}
