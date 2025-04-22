package adapter

import (
	"be-yourmoments/transaction-svc/internal/model"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentAdapter interface {
	CreateSnapshot(ctx context.Context, request *model.PaymentSnapshotRequest) (*snap.Response, error)
	GetPaymentServerKey() string
}

type paymentAdapter struct {
	snapClient *snap.Client
}

func NewPaymentAdapter(snapClient *snap.Client) PaymentAdapter {
	return &paymentAdapter{
		snapClient: snapClient,
	}
}

func (a *paymentAdapter) CreateSnapshot(ctx context.Context, request *model.PaymentSnapshotRequest) (*snap.Response, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  request.OrderID,
			GrossAmt: request.GrossAmount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			Email: request.Email,
		},
	}

	type result struct {
		resp *snap.Response
		err  error
	}

	resultChan := make(chan *result, 1)

	go func() {
		resp, err := a.snapClient.CreateTransaction(snapReq)
		if err != nil {
			log.Printf("error happenss: %v", err)
			resultChan <- &result{resp: nil, err: err}
			return
		}
		resultChan <- &result{resp: resp, err: nil}
	}()

	select {
	case <-timeoutCtx.Done():
		return nil, fmt.Errorf("create snapshot timeouts")
	case res := <-resultChan:
		if res.err != nil {
			log.Print("error happenss")
			return nil, fmt.Errorf("create snapshot errors: %w", res.err)
		}
		return res.resp, nil
	}

}

func (a *paymentAdapter) GetPaymentServerKey() string {
	return a.snapClient.ServerKey
}
