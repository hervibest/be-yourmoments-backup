package scheduler

import (
	"context"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/usecase"
)

type SchedulerRunner interface {
	Start()
}

type schedulerRunner struct {
	scheduler gocron.Scheduler
	usecase   usecase.SchedulerUseCase
	logs      *logger.Log
}

func NewSchedulerRunner(s gocron.Scheduler, usecase usecase.SchedulerUseCase, logs *logger.Log) SchedulerRunner {
	return &schedulerRunner{scheduler: s, usecase: usecase, logs: logs}
}

func (r *schedulerRunner) Start() {
	// 1. Define a “run every 5 minutes” job
	jobDef := gocron.DurationJob(5 * time.Minute)

	// 2. Register it with the scheduler
	_, err := r.scheduler.NewJob(
		jobDef,
		gocron.NewTask(func(ctx context.Context) {
			// give the task up to 4 minutes before timing out
			ctx, cancel := context.WithTimeout(ctx, 4*time.Minute)
			defer cancel()

			if err := r.usecase.CheckTransactionStatus(ctx); err != nil {
				// TODO: replace with your logger
				r.logs.CustomError("Scheduler error:", err.Error())
			}
		}),
	)
	if err != nil {
		r.logs.CustomError("Scheduler error:", err.Error())
	}

	r.logs.Log("Starting scheduler for every 5 minute:")

	// 3. Start the scheduler (asynchronously)
	r.scheduler.Start()
}
