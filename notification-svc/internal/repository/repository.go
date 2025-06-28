package repository

import (
	"context"
	"database/sql"

	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper"
	errorcode "github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper/enum/error"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper/logger"
	"github.com/jmoiron/sqlx"
)

type Querier interface {
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

type BeginTx interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (TransactionTx, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

type databaseAdapter struct {
	db *sqlx.DB
}

func NewDatabaseAdapter(db *sqlx.DB) BeginTx {
	return &databaseAdapter{
		db: db,
	}
}

func (r *databaseAdapter) BeginTxx(ctx context.Context, opts *sql.TxOptions) (TransactionTx, error) {
	return r.db.BeginTxx(ctx, opts)
}

func (r *databaseAdapter) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return r.db.QueryRowxContext(ctx, query, args...)
}

func (r *databaseAdapter) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return r.db.QueryxContext(ctx, query, args...)
}

func (r *databaseAdapter) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return r.db.GetContext(ctx, dest, query, args...)
}

func (r *databaseAdapter) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return r.db.SelectContext(ctx, dest, query, args...)
}

func (r *databaseAdapter) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return r.db.ExecContext(ctx, query, args...)
}

type TransactionTx interface {
	Commit() error
	Rollback() error
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func BeginTransaction(ctx context.Context, logs logger.Log, db BeginTx, fn func(tx TransactionTx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		// logs.Error(fmt.Sprintf("failed begin transaction %s", err.Error()), &logger.Options{
		// 	IsPrintStack: true,
		// })
		return helper.NewAppError(errorcode.ErrInternal, "Something went wrong. Please try again later", err)
	}

	rolledBack := false
	defer func() {
		if err != nil && !rolledBack {
			_ = tx.Rollback()
		}
	}()

	if err = fn(tx); err != nil {
		_ = tx.Rollback()
		rolledBack = true
		return err
	}

	if err = tx.Commit(); err != nil {
		// logs.Error(fmt.Sprintf("failed to commit transaction %s", err.Error()), &logger.Options{
		// 	IsPrintStack: true,
		// })
		return helper.NewAppError(errorcode.ErrInternal, "Something went wrong. Please try again later", err)
	}
	return nil
}
