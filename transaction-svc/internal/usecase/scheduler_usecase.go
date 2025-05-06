package usecase

import (
	"context"
	"fmt"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/repository"
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
	transactionUseCase    transactionUseCase
	logs                  *logger.Log
}

func NewSchedulerUseCase(db *sqlx.DB, transactionRepository repository.TransactionRepository, paymentAdapter adapter.PaymentAdapter, logs *logger.Log) SchedulerUseCase {
	return &schedulerUseCase{db: db, transactionRepository: transactionRepository, paymentAdapter: paymentAdapter, logs: logs}
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

	checkJobs := make(chan *entity.Transaction, 10)
	updateJobs := make(chan *Job, 10)
	var wgCheck, wgUpdate sync.WaitGroup

	// Stage 1: Worker pool untuk CheckTransactionStatus
	for i := 0; i < 10; i++ {
		wgCheck.Add(1)
		go func() {
			defer wgCheck.Done()
			for tx := range checkJobs {
				resp, err := u.paymentAdapter.CheckTransactionStatus(ctx, tx.Id)
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
				requestIsValid, hashedSignature := u.transactionUseCase.CheckPaymentSignature(
					job.response.SignatureKey,
					job.transaction.Id,
					job.response.StatusCode,
					job.response.GrossAmount,
				)

				if !requestIsValid {
					u.logs.Log(fmt.Sprintf("Invalid signature: expected=%s got=%s", hashedSignature, job.response.SignatureKey))
					continue
				}

				jsonValue, err := sonic.ConfigFastest.Marshal(job.response)
				if err != nil {
					u.logs.Log(fmt.Sprintf("marshal user : %+v", err))
					continue
				}

				_ = u.transactionUseCase.CheckAndUpdateTransaction(
					ctx, jsonValue, job.response.SettlementTime,
					job.response.TransactionStatus, job.transaction,
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
