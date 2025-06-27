package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/entity"
	errorcode "github.com/hervibest/be-yourmoments-backup/user-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper"
	mockadapter "github.com/hervibest/be-yourmoments-backup/user-svc/internal/mocks/adapter"
	mockproducer "github.com/hervibest/be-yourmoments-backup/user-svc/internal/mocks/gateway/producer"
	mocklogger "github.com/hervibest/be-yourmoments-backup/user-svc/internal/mocks/helper/logger"
	mockrepository "github.com/hervibest/be-yourmoments-backup/user-svc/internal/mocks/repository"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/usecase"
	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
	// Import package internal (sesuaikan path-nya)
)

func TestRegisterByPhoneNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	// Buat mock untuk dependency
	mockUserRepo := mockrepository.NewMockUserRepository(ctrl)
	mockUserProfileRepo := mockrepository.NewMockUserProfileRepository(ctrl)
	mockEmailVerificationRepo := mockrepository.NewMockEmailVerificationRepository(ctrl)
	mockResetPasswordRepo := mockrepository.NewMockResetPasswordRepository(ctrl)
	mockUserDeviceRepo := mockrepository.NewMockUserDeviceRepository(ctrl)

	mockDB := mockrepository.NewMockBeginTx(ctrl)
	mockTx := mockrepository.NewMockTransactionTx(ctrl)

	mockCacheAdapter := mockadapter.NewMockCacheAdapter(ctrl)
	mockEmailAdapter := mockadapter.NewMockEmailAdapter(ctrl)
	mockGoogleTokenAdapter := mockadapter.NewMockGoogleTokenAdapter(ctrl)
	mockJwtAdapter := mockadapter.NewMockJWTAdapter(ctrl)
	mockSecurityAdapter := mockadapter.NewMockSecurityAdapter(ctrl)
	mockRealtimeChatAdapter := mockadapter.NewMockRealtimeChatAdapter(ctrl)
	mockUserProducer := mockproducer.NewMockUserProducer(ctrl)

	mockLog := mocklogger.NewMockLog(ctrl)

	authUC := usecase.NewAuthUseCase(mockDB, mockUserRepo, mockUserProfileRepo, mockEmailVerificationRepo, mockResetPasswordRepo,
		mockUserDeviceRepo, mockGoogleTokenAdapter, mockEmailAdapter, mockJwtAdapter, mockSecurityAdapter, mockCacheAdapter,
		mockRealtimeChatAdapter, mockUserProducer, mockLog)
	// Data request testing
	now := time.Now()
	req := &model.RegisterByPhoneRequest{
		PhoneNumber: "08123456789",
		Username:    "testuser",
		Password:    "secret",
		BirthDate:   &now,
	}

	t.Run("Phone number already taken", func(t *testing.T) {
		// Ekspektasi: jumlah data berdasarkan phone number > 0
		mockUserRepo.EXPECT().CountByPhoneNumber(ctx, req.PhoneNumber).Return(1, nil)

		resp, err := authUC.RegisterByPhoneNumber(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)

		// Pastikan error-nya sesuai dengan fiber.NewError(http.StatusBadRequest, "phone number has already taken")
		appErr, ok := err.(*helper.AppError)
		assert.True(t, ok)
		assert.Equal(t, errorcode.ErrAlreadyExists, appErr.Code)
		assert.Equal(t, "Phone number has already been taken", appErr.Message)
	})

	t.Run("Username already taken", func(t *testing.T) {
		// Nomor telepon tidak terdaftar
		mockUserRepo.EXPECT().CountByPhoneNumber(ctx, req.PhoneNumber).Return(0, nil)
		// Ekspektasi: username sudah ada
		mockUserRepo.EXPECT().CountByUsername(ctx, req.Username).Return(1, nil)

		resp, err := authUC.RegisterByPhoneNumber(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)

		appErr, ok := err.(*helper.AppError)
		assert.True(t, ok)
		assert.Equal(t, errorcode.ErrAlreadyExists, appErr.Code)
		assert.Equal(t, "Username has already been taken", appErr.Message)
	})

	t.Run("Successful registration", func(t *testing.T) {
		// Nomor telepon dan username belum terdaftar
		mockUserRepo.EXPECT().CountByPhoneNumber(ctx, req.PhoneNumber).Return(0, nil)
		mockUserRepo.EXPECT().CountByUsername(ctx, req.Username).Return(0, nil)

		// Mulai transaksi
		mockDB.EXPECT().BeginTxx(ctx, gomock.Any()).Return(mockTx, nil)

		// Ekspektasi pembuatan user; gunakan matcher Any() untuk parameter user
		var createdUser *entity.User
		mockUserRepo.EXPECT().CreateByPhoneNumber(ctx, mockTx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, tx interface{}, user *entity.User) (*entity.User, error) {
				createdUser = user
				// Simulasikan kembalian user yang telah dibuat
				return user, nil
			},
		)

		// Ekspektasi pembuatan user profile
		mockUserProfileRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, tx interface{}, profile *entity.UserProfile) (*entity.UserProfile, error) {
				// Misal set id profile dan pastikan UserId sesuai dengan user yang dibuat
				profile.Id = "profile-123"
				profile.UserId = createdUser.Id
				return profile, nil
			},
		)

		// Ekspektasi commit transaksi
		mockTx.EXPECT().Commit().Return(nil)

		resp, err := authUC.RegisterByPhoneNumber(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		// Verifikasi field yang diharapkan, misalnya username
		assert.Equal(t, req.Username, resp.Username)
		// Asersi tambahan dapat disesuaikan dengan hasil konversi dari converter.UserToResponse
	})

	t.Run("Transaction internal commit error", func(t *testing.T) {
		// Skenario: commit transaksi gagal

		// Nomor telepon dan username belum terdaftar
		mockUserRepo.EXPECT().CountByPhoneNumber(ctx, req.PhoneNumber).Return(0, nil)
		mockUserRepo.EXPECT().CountByUsername(ctx, req.Username).Return(0, nil)

		// Mulai transaksi
		mockDB.EXPECT().BeginTxx(ctx, gomock.Any()).Return(mockTx, nil)

		// Ekspektasi pembuatan user
		mockUserRepo.EXPECT().CreateByPhoneNumber(ctx, mockTx, gomock.Any()).Return(&entity.User{
			Id:          "user-123",
			Username:    req.Username,
			Password:    sql.NullString{String: "hashed"},
			PhoneNumber: sql.NullString{String: req.PhoneNumber},
			CreatedAt:   &time.Time{},
			UpdatedAt:   &time.Time{},
		}, nil)

		// Ekspektasi pembuatan user profile
		mockUserProfileRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).Return(&entity.UserProfile{
			Id: "profile-123",
		}, nil)

		// Simulasikan error pada commit transaksi
		commitErr := fmt.Errorf("commit failed")
		mockTx.EXPECT().Commit().Return(commitErr)
		mockTx.EXPECT().Rollback().Return(nil)

		resp, err := authUC.RegisterByPhoneNumber(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		appErr, ok := err.(*helper.AppError)
		assert.True(t, ok)
		assert.Equal(t, errorcode.ErrInternal, appErr.Code)
		assert.Equal(t, "Something went wrong. Please try again later", appErr.Message)
	})

	t.Run("Fn failed transaction commit error", func(t *testing.T) {
		mockUserRepo.EXPECT().CountByPhoneNumber(ctx, req.PhoneNumber).Return(0, nil)
		mockUserRepo.EXPECT().CountByUsername(ctx, req.Username).Return(0, nil)

		mockDB.EXPECT().BeginTxx(ctx, gomock.Any()).Return(mockTx, nil)

		mockUserRepo.EXPECT().CreateByPhoneNumber(ctx, mockTx, gomock.Any()).Return(&entity.User{
			Id:          "user-123",
			Username:    req.Username,
			Password:    sql.NullString{String: "hashed"},
			PhoneNumber: sql.NullString{String: req.PhoneNumber},
			CreatedAt:   &time.Time{},
			UpdatedAt:   &time.Time{},
		}, nil)

		txErr := errors.New("tx operation failed")

		mockUserProfileRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).Return(nil, txErr)
		mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

		// Simulasikan error pada commit transaksi
		mockTx.EXPECT().Rollback().Return(nil)

		resp, err := authUC.RegisterByPhoneNumber(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		appErr, ok := err.(*helper.AppError)
		assert.True(t, ok)
		assert.Equal(t, errorcode.ErrInternal, appErr.Code)
		assert.Equal(t, "Something went wrong. Please try again later", appErr.Message)
	})
}

func TestRegisterByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	// Buat mock untuk dependency
	mockUserRepo := mockrepository.NewMockUserRepository(ctrl)
	mockUserProfileRepo := mockrepository.NewMockUserProfileRepository(ctrl)
	mockEmailVerificationRepo := mockrepository.NewMockEmailVerificationRepository(ctrl)
	mockResetPasswordRepo := mockrepository.NewMockResetPasswordRepository(ctrl)
	mockUserDeviceRepo := mockrepository.NewMockUserDeviceRepository(ctrl)

	mockDB := mockrepository.NewMockBeginTx(ctrl)
	mockTx := mockrepository.NewMockTransactionTx(ctrl)

	mockCacheAdapter := mockadapter.NewMockCacheAdapter(ctrl)
	mockEmailAdapter := mockadapter.NewMockEmailAdapter(ctrl)
	mockGoogleTokenAdapter := mockadapter.NewMockGoogleTokenAdapter(ctrl)
	mockJwtAdapter := mockadapter.NewMockJWTAdapter(ctrl)
	mockSecurityAdapter := mockadapter.NewMockSecurityAdapter(ctrl)
	mockRealtimeChatAdapter := mockadapter.NewMockRealtimeChatAdapter(ctrl)
	mockUserProducer := mockproducer.NewMockUserProducer(ctrl)

	mockLog := mocklogger.NewMockLog(ctrl)

	authUC := usecase.NewAuthUseCase(mockDB, mockUserRepo, mockUserProfileRepo, mockEmailVerificationRepo, mockResetPasswordRepo,
		mockUserDeviceRepo, mockGoogleTokenAdapter, mockEmailAdapter, mockJwtAdapter, mockSecurityAdapter, mockCacheAdapter,
		mockRealtimeChatAdapter, mockUserProducer, mockLog)
	// Data request testing
	now := time.Now()
	req := &model.RegisterByEmailRequest{
		Email:     "hervipro@gmail.com",
		Username:  "testuser",
		Password:  "secret",
		BirthDate: &now,
	}

	t.Run("Email already taken", func(t *testing.T) {
		// Ekspektasi: jumlah data berdasarkan phone number > 0
		mockUserRepo.EXPECT().CountByEmail(ctx, req.Email).Return(1, nil)

		resp, err := authUC.RegisterByEmail(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)

		// Pastikan error-nya sesuai dengan fiber.NewError(http.StatusBadRequest, "phone number has already taken")
		appErr, ok := err.(*helper.AppError)
		assert.True(t, ok)
		assert.Equal(t, errorcode.ErrAlreadyExists, appErr.Code)
		assert.Equal(t, "Email has already been taken", appErr.Message)
	})

	t.Run("Username already taken", func(t *testing.T) {
		// Nomor telepon tidak terdaftar
		mockUserRepo.EXPECT().CountByEmail(ctx, req.Email).Return(0, nil)
		// Ekspektasi: username sudah ada
		mockUserRepo.EXPECT().CountByUsername(ctx, req.Username).Return(1, nil)

		resp, err := authUC.RegisterByEmail(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)

		appErr, ok := err.(*helper.AppError)
		assert.True(t, ok)
		assert.Equal(t, errorcode.ErrAlreadyExists, appErr.Code)
		assert.Equal(t, "Username has already been taken", appErr.Message)
	})

	t.Run("Commit failed", func(t *testing.T) {
		// Nomor telepon dan username belum terdaftar
		mockUserRepo.EXPECT().CountByEmail(ctx, req.Email).Return(0, nil)
		mockUserRepo.EXPECT().CountByUsername(ctx, req.Username).Return(0, nil)

		// Mulai transaksi
		mockDB.EXPECT().BeginTxx(ctx, gomock.Any()).Return(mockTx, nil)

		// Ekspektasi pembuatan user; gunakan matcher Any() untuk parameter user
		var createdUser *entity.User
		mockUserRepo.EXPECT().CreateByEmail(ctx, mockTx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, tx interface{}, user *entity.User) (*entity.User, error) {
				createdUser = user
				return user, nil
			},
		)

		mockUserProfileRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, tx interface{}, profile *entity.UserProfile) (*entity.UserProfile, error) {
				profile.Id = "profile-123"
				profile.UserId = createdUser.Id
				return profile, nil
			},
		)

		commitErr := fmt.Errorf("commit failed")
		mockTx.EXPECT().Commit().Return(commitErr)
		mockTx.EXPECT().Rollback().Return(nil)

		resp, err := authUC.RegisterByEmail(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		appErr, ok := err.(*helper.AppError)
		assert.True(t, ok)
		assert.Equal(t, errorcode.ErrInternal, appErr.Code)
		assert.Equal(t, "Something went wrong. Please try again later", appErr.Message)
	})

	t.Run("Successful registration", func(t *testing.T) {
		// Nomor telepon dan username belum terdaftar
		mockUserRepo.EXPECT().CountByEmail(ctx, req.Email).Return(0, nil)
		mockUserRepo.EXPECT().CountByUsername(ctx, req.Username).Return(0, nil)

		// Mulai transaksi
		mockDB.EXPECT().BeginTxx(ctx, gomock.Any()).Return(mockTx, nil)

		// Ekspektasi pembuatan user; gunakan matcher Any() untuk parameter user
		var createdUser *entity.User
		mockUserRepo.EXPECT().CreateByEmail(ctx, mockTx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, tx interface{}, user *entity.User) (*entity.User, error) {
				createdUser = user
				return user, nil
			},
		)

		mockUserProfileRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, tx interface{}, profile *entity.UserProfile) (*entity.UserProfile, error) {
				profile.Id = "profile-123"
				profile.UserId = createdUser.Id
				return profile, nil
			},
		)

		mockTx.EXPECT().Commit().Return(nil)

		mockEmailVerificationRepo.EXPECT().Insert(ctx, mockTx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, tx interface{}, emailRequest *entity.EmailVerification) (*entity.EmailVerification, error) {
				return nil, nil
			},
		)

		mockDB.EXPECT().BeginTxx(ctx, gomock.Any()).Return(mockTx, nil)
		mockTx.EXPECT().Commit().Return(nil)
		mockSecurityAdapter.EXPECT().Encrypt(gomock.Any()).Return("ENCRYPTED_PASSWORD", nil)
		mockEmailAdapter.EXPECT().SendEmail(req.Email, "ENCRYPTED_PASSWORD", "new email verification").Return(nil)

		resp, err := authUC.RegisterByEmail(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		// Verifikasi field yang diharapkan, misalnya username
		assert.Equal(t, req.Username, resp.Username)
		// Asersi tambahan dapat disesuaikan dengan hasil konversi dari converter.UserToResponse
	})

	// t.Run("Transaction internal commit error", func(t *testing.T) {
	// 	// Skenario: commit transaksi gagal

	// 	// Nomor telepon dan username belum terdaftar
	// 	mockUserRepo.EXPECT().CountByPhoneNumber(ctx, req.PhoneNumber).Return(0, nil)
	// 	mockUserRepo.EXPECT().CountByUsername(ctx, req.Username).Return(0, nil)

	// 	// Mulai transaksi
	// 	mockDB.EXPECT().BeginTxx(ctx, gomock.Any()).Return(mockTx, nil)

	// 	// Ekspektasi pembuatan user
	// 	mockUserRepo.EXPECT().CreateByPhoneNumber(ctx, mockTx, gomock.Any()).Return(&entity.User{
	// 		Id:          "user-123",
	// 		Username:    req.Username,
	// 		Password:    sql.NullString{String: "hashed"},
	// 		PhoneNumber: sql.NullString{String: req.PhoneNumber},
	// 		CreatedAt:   &time.Time{},
	// 		UpdatedAt:   &time.Time{},
	// 	}, nil)

	// 	// Ekspektasi pembuatan user profile
	// 	mockUserProfileRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).Return(&entity.UserProfile{
	// 		Id: "profile-123",
	// 	}, nil)

	// 	// Simulasikan error pada commit transaksi
	// 	commitErr := fmt.Errorf("commit failed")
	// 	mockTx.EXPECT().Commit().Return(commitErr)
	// 	mockTx.EXPECT().Rollback().Return(nil)

	// 	resp, err := authUC.RegisterByPhoneNumber(ctx, req)
	// 	assert.Error(t, err)
	// 	assert.Nil(t, resp)
	// 	appErr, ok := err.(*helper.AppError)
	// 	assert.True(t, ok)
	// 	assert.Equal(t, errorcode.ErrInternal, appErr.Code)
	// 	assert.Equal(t, "Something went wrong. Please try again later", appErr.Message)
	// })

	// t.Run("Fn failed transaction commit error", func(t *testing.T) {
	// 	mockUserRepo.EXPECT().CountByPhoneNumber(ctx, req.PhoneNumber).Return(0, nil)
	// 	mockUserRepo.EXPECT().CountByUsername(ctx, req.Username).Return(0, nil)

	// 	mockDB.EXPECT().BeginTxx(ctx, gomock.Any()).Return(mockTx, nil)

	// 	mockUserRepo.EXPECT().CreateByPhoneNumber(ctx, mockTx, gomock.Any()).Return(&entity.User{
	// 		Id:          "user-123",
	// 		Username:    req.Username,
	// 		Password:    sql.NullString{String: "hashed"},
	// 		PhoneNumber: sql.NullString{String: req.PhoneNumber},
	// 		CreatedAt:   &time.Time{},
	// 		UpdatedAt:   &time.Time{},
	// 	}, nil)

	// 	txErr := errors.New("tx operation failed")

	// 	mockUserProfileRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).Return(nil, txErr)
	// 	mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	// 	// Simulasikan error pada commit transaksi
	// 	mockTx.EXPECT().Rollback().Return(nil)

	// 	resp, err := authUC.RegisterByPhoneNumber(ctx, req)
	// 	assert.Error(t, err)
	// 	assert.Nil(t, resp)
	// 	appErr, ok := err.(*helper.AppError)
	// 	assert.True(t, ok)
	// 	assert.Equal(t, errorcode.ErrInternal, appErr.Code)
	// 	assert.Equal(t, "Something went wrong. Please try again later", appErr.Message)
	// })
}
