package adapter

import (
	"context"
	"log"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/discovery"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/utils"

	transcationpb "github.com/hervibest/be-yourmoments-backup/pb/transaction"
)

type TransactionAdapter interface {
	CreateWallet(ctx context.Context, creatorId string) (*transcationpb.CreateWalletResponse, error)
}

type transactionAdapter struct {
	client transcationpb.TransactionServiceClient
}

func NewTransactionAdapter(ctx context.Context, registry discovery.Registry, logs *logger.Log) (TransactionAdapter, error) {
	transactionServiceName := utils.GetEnv("TRANSACTION_SVC_NAME")
	conn, err := discovery.ServiceConnection(ctx, transactionServiceName, registry, logs)
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
