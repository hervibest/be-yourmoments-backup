package repository

import (
	"context"
	"errors"
	"testing"

	mockrepository "github.com/hervibest/be-yourmoments-backup/user-svc/internal/mocks/repository"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBeginTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockrepository.NewMockBeginTx(ctrl)
	mockTx := mockrepository.NewMockTransactionTx(ctrl)

	// Expectasi
	mockDB.EXPECT().
		BeginTxx(gomock.Any(), nil).
		Return(mockTx, nil)

	mockTx.EXPECT().Commit().Return(nil)

	// Test
	err := repository.BeginTransaction(context.Background(), nil, mockDB, func(tx repository.TransactionTx) error {
		return nil
	})

	assert.NoError(t, err)
}

func TestBeginTransaction_CommitFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockrepository.NewMockBeginTx(ctrl)
	mockTx := mockrepository.NewMockTransactionTx(ctrl)

	// Simulasi BeginTxx berhasil
	mockDB.EXPECT().
		BeginTxx(gomock.Any(), nil).
		Return(mockTx, nil)

	// Simulasi fn(tx) berhasil, tapi Commit gagal
	mockTx.EXPECT().Commit().Return(errors.New("commit failed"))
	mockTx.EXPECT().Rollback().Return(nil)

	// Optional: bisa juga expect error log (jika logger-nya bisa di-mock)

	// Test
	err := repository.BeginTransaction(context.Background(), nil, mockDB, func(tx repository.TransactionTx) error {
		return nil
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Something went wrong. Please try again later")
}

func TestBeginTransaction_FnFailed_RollbackCalled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockrepository.NewMockBeginTx(ctrl)
	mockTx := mockrepository.NewMockTransactionTx(ctrl)

	mockDB.EXPECT().
		BeginTxx(gomock.Any(), nil).
		Return(mockTx, nil)

	mockTx.EXPECT().Rollback().Return(nil)

	// Simulasi kesalahan dalam fungsi transaksi
	txErr := errors.New("tx operation failed")

	err := repository.BeginTransaction(context.Background(), nil, mockDB, func(tx repository.TransactionTx) error {
		return txErr
	})

	assert.Equal(t, txErr, err)
}
