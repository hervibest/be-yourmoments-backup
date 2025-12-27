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
	transactionSvcName := utils.GetEnv("TRANSACTION_SVC_NAME")
	conn, err := discovery.NewGrpcClient(transactionSvcName)
	if err != nil {
		logs.CustomError("failed to connect transaction service with error : ", err)
		return nil, err
	}

	log.Printf("successfuly connected to %s", transactionSvcName)
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
