package config

import (
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/utils"

	"github.com/midtrans/midtrans-go"

	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransClient struct {
	Snap    *snap.Client
	CoreApi *coreapi.Client
}

func NewMidtransClient() *MidtransClient {
	midtransKey := utils.GetEnv("MIDTRANS_SERVER_KEY")

	snap := &snap.Client{}
	snap.New(midtransKey, midtrans.Sandbox)

	coreApi := &coreapi.Client{}
	coreApi.New(midtransKey, midtrans.Sandbox)

	return &MidtransClient{
		Snap:    snap,
		CoreApi: coreApi,
	}
}
