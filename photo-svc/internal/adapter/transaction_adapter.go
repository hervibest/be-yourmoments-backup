package adapter

import (
	"be-yourmoments/photo-svc/internal/helper"
	"be-yourmoments/photo-svc/internal/helper/discovery"
	"context"
	"log"

	transcationpb "github.com/be-yourmoments/pb/transaction"
)

type TransactionAdapter interface {
	CreateWallet(ctx context.Context, creatorId string) (*transcationpb.CreateWalletResponse, error)
}

type transactionAdapter struct {
	client transcationpb.TransactionServiceClient
}

func NewTransactionAdapter(ctx context.Context, registry discovery.Registry) (TransactionAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "transaction-svc-grpc", registry)
	if err != nil {
		return nil, err
	}

	log.Print("successfuly connected to transaction-svc-grpc")
	client := transcationpb.NewTransactionServiceClient(conn)

	return &transactionAdapter{
		client: client,
	}, nil
}

func (a *transactionAdapter) CreateWallet(ctx context.Context, creatorId string) (*transcationpb.CreateWalletResponse, error) {
	processPhotoRequest := &transcationpb.CreateWalletRequest{
		CreatorId: creatorId,
	}

	response, err := a.client.CreateWallet(ctx, processPhotoRequest)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	return response, nil
}
