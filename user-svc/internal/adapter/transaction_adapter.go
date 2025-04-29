package adapter

import (
	"context"
	"log"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper"
	discovery "github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/discovery"

	transcationpb "github.com/hervibest/be-yourmoments-backup/pb/transaction"
)

type TransactionAdapter interface {
	CreateWallet(ctx context.Context, creatorId string) (*entity.Wallet, error)
	GetWallet(ctx context.Context, creatorId string) (*entity.Wallet, error)
}

type transactionAdapter struct {
	client transcationpb.TransactionServiceClient
}

func NewTransactionAdapter(ctx context.Context, registry discovery.Registry) (TransactionAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "transaction-svc-grpc", registry)
	if err != nil {
		return nil, err
	}

	log.Print("transaction-svc-grpc")
	client := transcationpb.NewTransactionServiceClient(conn)

	return &transactionAdapter{
		client: client,
	}, nil
}

func (a *transactionAdapter) CreateWallet(ctx context.Context, creatorId string) (*entity.Wallet, error) {
	pbRequest := &transcationpb.CreateWalletRequest{
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
	pbRequest := &transcationpb.GetWalletRequest{
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
