package adapter

import (
	"be-yourmoments/photo-svc/internal/helper"
	"be-yourmoments/photo-svc/internal/helper/discovery"
	"context"
	"log"

	"github.com/be-yourmoments/pb"
)

type TransactionAdapter interface {
	CreateWallet(ctx context.Context, creatorId string) (*pb.CreateWalletResponse, error)
}

type transactionAdapter struct {
	client pb.TransactionServiceClient
}

func NewTransactionAdapter(ctx context.Context, registry discovery.Registry) (TransactionAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "transaction-svc-grpc", registry)
	if err != nil {
		return nil, err
	}

	log.Print("successfuly connected to transaction-svc-grpc")
	client := pb.NewTransactionServiceClient(conn)

	return &transactionAdapter{
		client: client,
	}, nil
}

func (a *transactionAdapter) CreateWallet(ctx context.Context, creatorId string) (*pb.CreateWalletResponse, error) {
	processPhotoRequest := &pb.CreateWalletRequest{
		CreatorId: creatorId,
	}

	response, err := a.client.CreateWallet(ctx, processPhotoRequest)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	return response, nil
}
