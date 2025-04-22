package grpc

import (
	"be-yourmoments/transaction-svc/internal/helper"
	"be-yourmoments/transaction-svc/internal/model"
	"be-yourmoments/transaction-svc/internal/usecase"
	"context"
	"log"
	"net/http"

	"github.com/be-yourmoments/pb"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TransactionGRPCHandler struct {
	walletUseCase usecase.WalletUsecase
	pb.UnimplementedTransactionServiceServer
}

func NewTransactionGRPCHandler(server *grpc.Server, walletUseCase usecase.WalletUsecase) {
	handler := &TransactionGRPCHandler{
		walletUseCase: walletUseCase,
	}

	pb.RegisterTransactionServiceServer(server, handler)

	if walletUseCase == nil {
		log.Print("wallet didnt initialized in constructor")
	}
}

func (h *TransactionGRPCHandler) CreateWallet(ctx context.Context, pbReq *pb.CreateWalletRequest) (
	*pb.CreateWalletResponse, error) {
	log.Println("----  CreatePhoto Requets via GRPC in transaction-svc ------")

	request := &model.RequestCreateWallet{
		CreatorId: pbReq.GetCreatorId(),
	}

	if h.walletUseCase == nil {
		log.Print("wallet didnt initialized in create wallet method")
	}

	log.Print(request.CreatorId)
	wallet, err := h.walletUseCase.CreateWallet(context.Background(), request)
	if err != nil {
		return nil, helper.ErrGRPC(err)
	}

	pbWallet := &pb.Wallet{
		Id:        wallet.Id,
		CreatorId: wallet.CreatorId,
		CreatedAt: &timestamppb.Timestamp{
			Seconds: wallet.CreatedAt.Unix(),
			Nanos:   int32(wallet.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: wallet.UpdatedAt.Unix(),
			Nanos:   int32(wallet.UpdatedAt.Nanosecond()),
		},
	}

	return &pb.CreateWalletResponse{
		Status: http.StatusCreated,
		Wallet: pbWallet,
	}, nil
}
