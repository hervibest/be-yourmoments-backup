package adapter

import (
	"be-yourmoments/user-svc/internal/entity"
	"be-yourmoments/user-svc/internal/helper"
	discovery "be-yourmoments/user-svc/internal/helper/discovery"
	"context"
	"log"

	"github.com/be-yourmoments/pb"
)

type TransactionAdapter interface {
	CreateWallet(ctx context.Context, creatorId string) (*entity.Wallet, error)
	GetWallet(ctx context.Context, creatorId string) (*entity.Wallet, error)
}

type transactionAdapter struct {
	client pb.TransactionServiceClient
}

func NewTransactionAdapter(ctx context.Context, registry discovery.Registry) (TransactionAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "transaction-svc-grpc", registry)
	if err != nil {
		return nil, err
	}

	log.Print("transaction-svc-grpc")
	client := pb.NewTransactionServiceClient(conn)

	return &transactionAdapter{
		client: client,
	}, nil
}

func (a *transactionAdapter) CreateWallet(ctx context.Context, creatorId string) (*entity.Wallet, error) {
	pbRequest := &pb.CreateWalletRequest{
		CreatorId: creatorId,
	}

	pbResponse, err := a.client.CreateWallet(context.Background(), pbRequest)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	wallet := &entity.Wallet{
		Id:        pbResponse.GetWallet().GetId(),
		CreatorId: pbResponse.GetWallet().GetCreatorId(),
		Balance:   pbResponse.GetWallet().GetBalance(),
	}

	return wallet, nil
}

func (a *transactionAdapter) GetWallet(ctx context.Context, creatorId string) (*entity.Wallet, error) {
	pbRequest := &pb.GetWalletRequest{
		CreatorId: creatorId,
	}

	log.Print("[user-svc] Get wallet by creator id :", creatorId)

	pbResponse, err := a.client.GetWallet(context.Background(), pbRequest)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	wallet := &entity.Wallet{
		Id:        pbResponse.GetWallet().GetId(),
		CreatorId: pbResponse.GetWallet().GetCreatorId(),
		Balance:   pbResponse.GetWallet().GetBalance(),
	}

	return wallet, nil
}
