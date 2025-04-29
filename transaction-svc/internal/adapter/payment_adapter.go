package adapter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sony/gobreaker"
)

type PaymentAdapter interface {
	CreateSnapshot(ctx context.Context, request *model.PaymentSnapshotRequest) (*snap.Response, error)
	GetPaymentServerKey() string
}

type paymentAdapter struct {
	snapClient     *snap.Client
	circuitBreaker *gobreaker.CircuitBreaker
	logs           *logger.Log
}

func NewPaymentAdapter(snapClient *snap.Client, logs *logger.Log) PaymentAdapter {
	cbSettings := gobreaker.Settings{
		Name:        "MidtransSnapshot",
		MaxRequests: 3,                // max concurrent request saat half-open
		Interval:    60 * time.Second, // reset counter tiap 1 menit
		Timeout:     10 * time.Second, // setelah trip, coba buka lagi setelah 10s
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.Requests >= 5 && float64(counts.TotalFailures)/float64(counts.Requests) >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("[CircuitBreaker] %s: %s -> %s\n", name, from.String(), to.String())
		},
	}

	cb := gobreaker.NewCircuitBreaker(cbSettings)

	return &paymentAdapter{
		snapClient:     snapClient,
		circuitBreaker: cb,
		logs:           logs,
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

	res, err := a.circuitBreaker.Execute(func() (interface{}, error) {
		resultChan := make(chan *snap.Response, 1)
		errChan := make(chan error, 1)

		go func() {
			resp, err := a.snapClient.CreateTransaction(snapReq)
			if err != nil {
				errChan <- err
				return
			}
			resultChan <- resp
		}()

		select {
		case <-timeoutCtx.Done():
			return nil, fmt.Errorf("midtrans request timeout: %w", timeoutCtx.Err())
		case err := <-errChan:
			return nil, err
		case resp := <-resultChan:
			return resp, nil
		}
	})

	if err != nil {
		a.logs.CustomError("[PaymentAdapter] CreateSnapshot error: %v", err)
		return nil, fmt.Errorf("midtrans create snapshot error: %w", err)
	}

	return res.(*snap.Response), nil
}

func (a *paymentAdapter) GetPaymentServerKey() string {
	return a.snapClient.ServerKey
}
