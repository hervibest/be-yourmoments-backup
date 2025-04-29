package config

import (
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/utils"

	"github.com/midtrans/midtrans-go"

	"github.com/midtrans/midtrans-go/snap"
)

func NewMidtransClient() *snap.Client {
	midtransKey := utils.GetEnv("MIDTRANS_SERVER_KEY")
	snapClient := &snap.Client{}
	snapClient.New(midtransKey, midtrans.Sandbox)

	return snapClient
}
