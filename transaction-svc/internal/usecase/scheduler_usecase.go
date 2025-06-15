package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/enum"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model/converter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/repository"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/usecase/contract"
	"github.com/jmoiron/sqlx"
	"github.com/midtrans/midtrans-go/coreapi"
)

type SchedulerUseCase interface {
	CheckTransactionStatus(ctx context.Context) error
}
type schedulerUseCase struct {
	db                    *sqlx.DB
	transactionRepository repository.TransactionRepository
	paymentAdapter        adapter.PaymentAdapter
	transactionUseCase    contract.TransactionUseCase
	cancelationUseCase    contract.CancelationUseCase
	logs                  *logger.Log
}

func NewSchedulerUseCase(db *sqlx.DB, transactionRepository repository.TransactionRepository, transactionUseCase contract.TransactionUseCase,
	cancelationUseCase contract.CancelationUseCase, paymentAdapter adapter.PaymentAdapter, logs *logger.Log) SchedulerUseCase {
	return &schedulerUseCase{db: db, transactionRepository: transactionRepository, transactionUseCase: transactionUseCase,
		cancelationUseCase: cancelationUseCase, paymentAdapter: paymentAdapter, logs: logs}
}

type Job struct {
	response    *coreapi.TransactionStatusResponse
	transaction *entity.Transaction
}

func (u *schedulerUseCase) CheckTransactionStatus(ctx context.Context) error {
	transactions, err := u.transactionRepository.FindManyCheckable(ctx, u.db)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to find many checkable transaction", err)
	}

	if len(*transactions) == 0 {
		u.logs.Log("There are no checkable transaction, ignoring check transaction process")
		return nil
	}

	checkJobs := make(chan *entity.Transaction, 10)
	updateJobs := make(chan *Job, 10)
	var wgCheck, wgUpdate sync.WaitGroup

	// Stage 1: Worker pool untuk CheckTransactionStatus
	for i := 0; i < 10; i++ {
		wgCheck.Add(1)
		go func() {
			defer wgCheck.Done()
			for tx := range checkJobs {
				if tx.CreatedAt != nil && tx.Status != enum.TransactionStatusExpired {
					if time.Since(*tx.CreatedAt) > 15*time.Minute {
						// Mark transaksi sebagai expired
						err := u.cancelationUseCase.ExpirePendingTransaction(ctx, tx.Id)
						if err != nil {
							log.Printf("Gagal meng-expire transaksi %s: %v", tx.Id, err)
						}
						continue // skip pengecekan ke payment gateway
					}
				}
				resp, err := u.paymentAdapter.CheckTransactionStatus(context.Background(), tx.Id)
				if err != nil {
					u.logs.Log(fmt.Sprintf("failed check status for %s: %v", tx.Id, err))
					continue
				}

				updateJobs <- &Job{
					response:    resp,
					transaction: tx,
				}
			}
		}()
	}

	// Stage 2: Worker pool untuk CheckAndUpdateTransaction
	for i := 0; i < 5; i++ {
		wgUpdate.Add(1)
		go func() {
			defer wgUpdate.Done()
			for job := range updateJobs {
				jsonValue, err := sonic.ConfigFastest.Marshal(job.response)
				if err != nil {
					u.logs.Log(fmt.Sprintf("marshal user : %+v", err))
					continue
				}

				request := converter.SchedulerReqToCheckAndUpdate(job.response, jsonValue)
				_ = u.transactionUseCase.CheckAndUpdateTransaction(
					ctx, request,
				)
			}
		}()
	}

	// Kirim job transaksi
	go func() {
		for _, tx := range *transactions {
			checkJobs <- tx
		}
		close(checkJobs)
	}()

	// Tutup updateJobs setelah semua check selesai
	go func() {
		wgCheck.Wait()
		close(updateJobs)
	}()

	wgUpdate.Wait()
	return nil
}
