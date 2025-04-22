package repository

import (
	errorcode "be-yourmoments/photo-svc/internal/enum/error"
	"be-yourmoments/photo-svc/internal/helper"
	"be-yourmoments/photo-svc/internal/helper/logger"
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Querier interface {
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error

	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

type BeginTx interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func BeginTxx(db BeginTx, ctx context.Context, logs *logger.Log) (*sqlx.Tx, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		logs.Error(fmt.Sprintf("failed begin transaction %s", err.Error()), &logger.Options{
			IsPrintStack: true,
		})
		return nil, helper.NewAppError(errorcode.ErrInternal, "Something went wrong. Please try again later", err)
	}
	return tx, nil
}

func Rollback(err error, tx *sqlx.Tx, ctx context.Context, logs *logger.Log) {
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logs.Error(fmt.Sprintf("failed to rollback transaction %s", err.Error()), &logger.Options{
				IsPrintStack: true,
			})
		}
	}
}

func Commit(tx *sqlx.Tx, logs *logger.Log) error {
	if err := tx.Commit(); err != nil {
		logs.Error(fmt.Sprintf("failed to commit transaction %s", err.Error()), &logger.Options{
			IsPrintStack: true,
		})
		return helper.NewAppError(errorcode.ErrInternal, "Something went wrong. Please try again later", err)
	}
	return nil
}
