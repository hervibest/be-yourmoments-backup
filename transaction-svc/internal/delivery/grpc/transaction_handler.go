package grpc

import (
	"context"
	"log"
	"net/http"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/usecase"

	transactionpb "github.com/hervibest/be-yourmoments-backup/pb/transaction"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TransactionGRPCHandler struct {
	walletUseCase usecase.WalletUsecase
	transactionpb.UnimplementedTransactionServiceServer
}

func NewTransactionGRPCHandler(server *grpc.Server, walletUseCase usecase.WalletUsecase) {
	handler := &TransactionGRPCHandler{
		walletUseCase: walletUseCase,
	}

	transactionpb.RegisterTransactionServiceServer(server, handler)

	if walletUseCase == nil {
		log.Print("wallet didnt initialized in constructor")
	}
}

func (h *TransactionGRPCHandler) CreateWallet(ctx context.Context, pbReq *transactionpb.CreateWalletRequest) (
	*transactionpb.CreateWalletResponse, error) {
	log.Println("----  Create Wallet GRPC in transaction-svc ------")

	request := &model.CreateWalletRequest{
		CreatorId: pbReq.GetCreatorId(),
	}

	wallet, err := h.walletUseCase.CreateWallet(context.Background(), request)
	if err != nil {
		return nil, helper.ErrGRPC(err)
	}

	pbWallet := &transactionpb.Wallet{
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

	return &transactionpb.CreateWalletResponse{
		Status: http.StatusCreated,
		Wallet: pbWallet,
	}, nil
}

func (h *TransactionGRPCHandler) GetWallet(ctx context.Context, pbReq *transactionpb.GetWalletRequest) (
	*transactionpb.GetWalletResponse, error) {
	log.Println("----  Get Wallet Requets via GRPC in transaction-svc ------")

	request := &model.GetWalletRequest{
		CreatorId: pbReq.GetCreatorId(),
	}

	wallet, err := h.walletUseCase.GetWallet(context.Background(), request)
	if err != nil {
		return nil, helper.ErrGRPC(err)
	}

	pbWallet := &transactionpb.Wallet{
		Id:        wallet.Id,
		CreatorId: wallet.CreatorId,
	}

	return &transactionpb.GetWalletResponse{
		Status: http.StatusCreated,
		Wallet: pbWallet,
	}, nil
}
