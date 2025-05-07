package repository

import (
	"context"
	"database/sql"
	"fmt"

	errorcode "github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"

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

	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	// NamedStmt(stmt *sqlx.NamedStmt) *sqlx.NamedStmt
	// NamedStmtContext(ctx context.Context, stmt *sqlx.NamedStmt) *sqlx.NamedStmt
	Rebind(query string) string
}

type BeginTx interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func BeginTxx(db *sqlx.DB, ctx context.Context, logs *logger.Log) (*sqlx.Tx, error) {
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
		if tx != nil {
			_ = tx.Rollback() // abaikan error rollback
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

func BeginTransaction(db *sqlx.DB, ctx context.Context, logs *logger.Log, fn func(tx *sqlx.Tx) error) error {
	tx, err := BeginTxx(db, ctx, logs)
	if err != nil {
		return err
	}

	// Pastikan rollback hanya dipanggil jika Commit belum dilakukan
	defer func() {
		if p := recover(); p != nil {
			Rollback(nil, tx, ctx, logs) // panic tetap di-rollback
			panic(p)                     // re-throw panic
		} else if err != nil {
			Rollback(err, tx, ctx, logs)
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	if err = Commit(tx, logs); err != nil {
		return err
	}

	return nil
}
