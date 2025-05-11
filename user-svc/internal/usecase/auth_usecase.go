package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/enum"
	errorcode "github.com/hervibest/be-yourmoments-backup/user-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/nullable"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model/converter"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/repository"
	"github.com/redis/go-redis/v9"

	"cloud.google.com/go/firestore"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	AccessTokenRequest(ctx context.Context, refreshToken string) (*model.UserResponse, *model.TokenResponse, error)
	CreateDeviceToken(ctx context.Context, request *model.DeviceRequest) error
	Current(ctx context.Context, email string) (*model.UserResponse, error)
	Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, *model.TokenResponse, error)
	Logout(ctx context.Context, request *model.LogoutUserRequest) (bool, error)
	RegisterByEmail(ctx context.Context, request *model.RegisterByEmailRequest) (*model.UserResponse, error)
	RegisterOrLoginByGoogle(ctx context.Context, request *model.RegisterByGoogleRequest) (*model.UserResponse, *model.TokenResponse, error)
	RegisterByPhoneNumber(ctx context.Context, request *model.RegisterByPhoneRequest) (*model.UserResponse, error)
	RequestResetPassword(ctx context.Context, email string) error
	ResendEmailVerification(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, request *model.ResetPasswordUserRequest) error
	ValidateResetPassword(ctx context.Context, request *model.ValidateResetTokenRequest) (bool, error)
	Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.AuthResponse, error)
	VerifyEmail(ctx context.Context, request *model.VerifyEmailUserRequest) error
}

type authUseCase struct {
	db                    repository.BeginTx
	userRepository        repository.UserRepository
	userProfileRepository repository.UserProfileRepository
	emailVerificationRepo repository.EmailVerificationRepository
	resetPasswordRepo     repository.ResetPasswordRepository
	userDeviceRepository  repository.UserDeviceRepository
	googleTokenAdapter    adapter.GoogleTokenAdapter
	emailAdapter          adapter.EmailAdapter
	securityAdapter       adapter.SecurityAdapter
	jwtAdapter            adapter.JWTAdapter
	cacheAdapter          adapter.CacheAdapter
	firestoreAdapter      adapter.FirestoreClientAdapter
	photoAdapter          adapter.PhotoAdapter
	transactionAdapter    adapter.TransactionAdapter

	logs logger.Log
}

func NewAuthUseCase(db repository.BeginTx, userRepository repository.UserRepository, userProfileRepository repository.UserProfileRepository,
	emailVerificationRepo repository.EmailVerificationRepository, resetPasswordRepo repository.ResetPasswordRepository,
	userDeviceRepository repository.UserDeviceRepository, googleTokenAdapter adapter.GoogleTokenAdapter,
	emailAdapter adapter.EmailAdapter, jwtAdapter adapter.JWTAdapter, securityAdapter adapter.SecurityAdapter,
	cacheAdapter adapter.CacheAdapter, firestoreAdapter adapter.FirestoreClientAdapter,
	photoAdapter adapter.PhotoAdapter, transactionAdapter adapter.TransactionAdapter, logs logger.Log) AuthUseCase {
	return &authUseCase{
		db:                    db,
		userRepository:        userRepository,
		userProfileRepository: userProfileRepository,
		emailVerificationRepo: emailVerificationRepo,
		resetPasswordRepo:     resetPasswordRepo,
		userDeviceRepository:  userDeviceRepository,
		googleTokenAdapter:    googleTokenAdapter,
		emailAdapter:          emailAdapter,
		securityAdapter:       securityAdapter,
		jwtAdapter:            jwtAdapter,
		cacheAdapter:          cacheAdapter,
		firestoreAdapter:      firestoreAdapter,
		photoAdapter:          photoAdapter,
		transactionAdapter:    transactionAdapter,
		logs:                  logs,
	}
}

func (u *authUseCase) RegisterByPhoneNumber(ctx context.Context, request *model.RegisterByPhoneRequest) (*model.UserResponse, error) {
	countByNumberTotal, err := u.userRepository.CountByPhoneNumber(ctx, request.PhoneNumber)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to send count user by phone number", err)
	}

	//TODO query check cukup sekali saja untuk checking username dan phonenumber
	if countByNumberTotal > 0 {
		return nil, helper.NewUseCaseError(errorcode.ErrAlreadyExists, "Phone number has already been taken")
	}

	countByUsernameTotal, err := u.userRepository.CountByUsername(ctx, request.Username)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to send count user by username", err)
	}

	if countByUsernameTotal > 0 {
		return nil, helper.NewUseCaseError(errorcode.ErrAlreadyExists, "Username has already been taken")
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return nil, err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	now := time.Now()
	user := &entity.User{
		Id:       ulid.Make().String(),
		Username: request.Username,
		Password: sql.NullString{
			Valid:  true,
			String: string(hashedPassword),
		},
		PhoneNumber: sql.NullString{
			Valid:  true,
			String: request.PhoneNumber,
		},
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	user, err = u.userRepository.CreateByPhoneNumber(ctx, tx, user)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to create user by phone number", err)
	}

	userProfile := &entity.UserProfile{
		Id:        ulid.Make().String(),
		UserId:    user.Id,
		BirthDate: request.BirthDate,
		Nickname:  helper.GenerateNickname(),
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	_, err = u.userProfileRepository.Create(ctx, tx, userProfile)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to create user profile", err)
	}

	//WHAT TO DO WHATS APP BUSINESS VERIF OTP
	if err := repository.Commit(tx, u.logs); err != nil {
		return nil, err
	}

	return converter.UserToResponse(user), nil
}

// TODO EFICIENT QUERY ISSUE (COUNT COUNT AND COUNT)
func (u *authUseCase) RegisterOrLoginByGoogle(ctx context.Context, request *model.RegisterByGoogleRequest) (*model.UserResponse, *model.TokenResponse, error) {
	if request.Platform != enum.PlatformTypeWeb &&
		request.Platform != enum.PlatformTypeIOS &&
		request.Platform != enum.PlatformTypeAndroid {
		return nil, nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid Platform Type")
	}

	claims, err := u.googleTokenAdapter.ValidateGoogleToken(ctx, request.Token)
	if err != nil {
		return nil, nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, err.Error())
	}

	countByNotGoogleTotal, err := u.userRepository.CountByEmailNotGoogle(ctx, claims.Email)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to count by email not google", err)
	}

	if countByNotGoogleTotal > 0 {
		return nil, nil, helper.NewUseCaseError(errorcode.ErrAlreadyExists, "Email has already been takens")
	}

	countByGoogleTotal, err := u.userRepository.CountByEmailGoogleId(ctx, claims.Email, claims.GoogleId)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to count by email google id", err)
	}

	var user *entity.User
	var wallet *entity.Wallet
	var creator *entity.Creator
	var userProfile *entity.UserProfile
	var wg sync.WaitGroup

	// If user already registered
	if countByGoogleTotal > 0 {
		user, err = u.userRepository.FindByEmail(ctx, claims.Email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid email")
			}
			return nil, nil, helper.WrapInternalServerError(u.logs, "failed find user by email not google", err)
		}

		userProfRes, creatorRes, walletRes, err := u.getUseProfileWalletAndRepo(ctx, &wg, user.Id)
		if err != nil {
			return nil, nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid email")
		}

		userProfile = userProfRes
		wallet = walletRes
		creator = creatorRes

	} else {
		tx, err := repository.BeginTxx(u.db, ctx, u.logs)
		if err != nil {
			return nil, nil, err
		}

		defer func() {
			repository.Rollback(err, tx, ctx, u.logs)
		}()

		now := time.Now()
		user = &entity.User{
			Id: ulid.Make().String(),
			Email: sql.NullString{
				Valid:  true,
				String: claims.Email,
			},
			Username: claims.Username,
			GoogleId: sql.NullString{
				Valid:  true,
				String: claims.GoogleId,
			},
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		user, err = u.userRepository.CreateByGoogleSignIn(ctx, tx, user)
		if err != nil {
			return nil, nil, helper.WrapInternalServerError(u.logs, "failed to create user by google sign in", err)
		}

		userProfile = &entity.UserProfile{
			Id:     ulid.Make().String(),
			UserId: user.Id,
			// BirthDate: request.BirthDate,
			Nickname: helper.GenerateNickname(),
			ProfileUrl: sql.NullString{
				Valid:  true,
				String: claims.ProfilePictureUrl,
			},
			Similarity: uint(enum.DefaultSimilarityLevel),
			CreatedAt:  &now,
			UpdatedAt:  &now,
		}

		_, err = u.userProfileRepository.CreateWithProfileUrl(ctx, tx, userProfile)
		if err != nil {
			return nil, nil, helper.WrapInternalServerError(u.logs, "failed to create with profile url", err)
		}

		if err := repository.Commit(tx, u.logs); err != nil {
			return nil, nil, err
		}

		go func() {
			u.createChatRoom(context.Background(), user, userProfile)
		}()

		creator, wallet, err = u.createCreatorAndWallet(ctx, user.Id)
		if err != nil {
			return nil, nil, err
		}

		wg.Wait()
	}

	now := time.Now()
	userDevice := &entity.UserDevice{
		Id:        ulid.Make().String(),
		UserId:    user.Id,
		Token:     request.DeviceToken,
		Platform:  request.Platform,
		CreatedAt: &now,
	}

	_, err = u.userDeviceRepository.Create(ctx, u.db, userDevice)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to create user device", err)
	}

	setKey := fmt.Sprintf("fcm_tokens:%s", user.Id)
	if err := u.cacheAdapter.SAdd(ctx, setKey, request.DeviceToken); err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to SAdd redis set", err)
	}

	auth := &entity.Auth{
		Id:          user.Id,
		Username:    user.Username,
		Email:       user.Email.String,
		PhoneNumber: user.PhoneNumber.String,
		Similarity:  userProfile.Similarity,
		CreatorId:   creator.Id,
		WalletId:    wallet.Id,
	}

	token, err := u.generateToken(ctx, auth)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to generate token :", err)
	}

	return converter.UserToResponse(user), token, nil
}

func (u *authUseCase) getUseProfileWalletAndRepo(ctx context.Context, wg *sync.WaitGroup, userId string) (*entity.UserProfile, *entity.Creator, *entity.Wallet, error) {
	var (
		userProfile *entity.UserProfile
		creator     *entity.Creator
		wallet      *entity.Wallet
		err         error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		profile, err1 := u.userProfileRepository.FindByUserId(ctx, userId)
		if err1 != nil {
			if errors.Is(err1, sql.ErrNoRows) {
				u.logs.Error(fmt.Sprint("invalid user id for finding user profile : ", err1))
			}
			u.logs.Error(fmt.Sprint("failed find user by email not google : ", err1))
			return
		}
		userProfile = profile
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err1 error
		creator, err1 = u.photoAdapter.GetCreator(ctx, userId)
		if err1 != nil {
			u.logs.Error(fmt.Sprint("failed to get creator from photo service with error : ", err1))
			return
		}

		wallet, err1 = u.transactionAdapter.GetWallet(ctx, creator.Id)
		if err1 != nil {
			u.logs.Error(fmt.Sprint("failed to get wallet by creator id from wallet service with error : ", err1))
		}
	}()

	wg.Wait()
	return userProfile, creator, wallet, err
}

func (u *authUseCase) createChatRoom(ctx context.Context, user *entity.User, userProfile *entity.UserProfile) {
	//TODO apakah bisa dirapikan atau diwrap ke dalam adapter ?
	userRef := u.firestoreAdapter.Collection("users").Doc(user.Id)
	_, err := userRef.Get(ctx)
	if err != nil {
		_, err := userRef.Set(ctx, map[string]interface{}{
			"userId":     user.Id,
			"profileId":  userProfile.Id,
			"nickname":   userProfile.Nickname,
			"profileUrl": nullable.SQLStringToPtr(userProfile.ProfileUrl),
			"createdAt":  firestore.ServerTimestamp,
			"updatedAt":  firestore.ServerTimestamp,
		})
		if err != nil {
			u.logs.Error(fmt.Sprintf("Failed to create or get rooms from firebase when create user by google with err : %v and user id : %s", err, user.Id))
		}
	}
}

func (u *authUseCase) createCreatorAndWallet(ctx context.Context, userId string) (*entity.Creator, *entity.Wallet, error) {
	creator, err := u.photoAdapter.CreateCreator(ctx, userId)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to create creator", err)
	}

	wallet, err := u.transactionAdapter.CreateWallet(ctx, creator.Id)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to create wallet by creator id", err)
	}

	return creator, wallet, nil
}

func (u *authUseCase) RegisterByEmail(ctx context.Context, request *model.RegisterByEmailRequest) (*model.UserResponse, error) {
	countByNumberTotal, err := u.userRepository.CountByEmail(ctx, request.Email)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to count user by email", err)
	}

	if countByNumberTotal > 0 {
		return nil, helper.NewUseCaseError(errorcode.ErrAlreadyExists, "Email has already been taken")
	}

	countByUsernameTotal, err := u.userRepository.CountByUsername(ctx, request.Username)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to count user by username", err)
	}

	if countByUsernameTotal > 0 {
		return nil, helper.NewUseCaseError(errorcode.ErrAlreadyExists, "Username has already been takens")
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return nil, err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to generate hashed bcrypt", err)
	}

	now := time.Now()
	user := &entity.User{
		Id:       ulid.Make().String(),
		Username: request.Username,
		Email: sql.NullString{
			Valid:  true,
			String: request.Email,
		},
		Password: sql.NullString{
			Valid:  true,
			String: string(hashedPassword),
		},
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	user, err = u.userRepository.CreateByEmail(ctx, tx, user)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to create user by email verification :", err)
	}

	userProfile := &entity.UserProfile{
		Id:        ulid.Make().String(),
		UserId:    user.Id,
		BirthDate: request.BirthDate,
		Nickname:  helper.GenerateNickname(),
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	_, err = u.userProfileRepository.Create(ctx, tx, userProfile)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to send email verification :", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return nil, err
	}

	if err := u.requestEmailVerification(ctx, request.Email, true); err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to send email verification :", err)
	}

	return converter.UserToResponse(user), nil
}

func (u *authUseCase) requestEmailVerification(ctx context.Context, email string, newUser bool) error {
	token := uuid.NewString()
	now := time.Now()

	emailVerification := &entity.EmailVerification{
		Email:     email,
		Token:     token,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return nil
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if newUser {
		_, err = u.emailVerificationRepo.Insert(ctx, tx, emailVerification)
		if err != nil {
			return fmt.Errorf("insert email verification token : %+v", err)
		}
	} else {
		_, err = u.emailVerificationRepo.Update(ctx, tx, emailVerification)
		if err != nil {
			return fmt.Errorf("update email verification token : %+v", err)
		}
	}

	encryptedToken, err := u.securityAdapter.Encrypt(token)
	if err != nil {
		return fmt.Errorf("encrypt email verification token : %+v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction : %+v", err)
	}

	if err := u.emailAdapter.SendEmail(email, encryptedToken, "new email verification"); err != nil {
		return fmt.Errorf("send new email verification : %+v", err)
	}

	return nil
}

func (u *authUseCase) ResendEmailVerification(ctx context.Context, email string) error {
	user, err := u.userRepository.FindByEmailNotGoogle(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid email or password")
		}
		return helper.WrapInternalServerError(u.logs, "failed to find user by email :", err)
	}

	if user.HasVerifiedEmail() {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "User already verified")
	}

	_, err = u.emailVerificationRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid email or password")
		}
		return helper.WrapInternalServerError(u.logs, "failed to find email verification by email", err)
	}

	if err := u.requestEmailVerification(ctx, email, false); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to send email verification", err)
	}

	return nil
}

func (u *authUseCase) VerifyEmail(ctx context.Context, request *model.VerifyEmailUserRequest) error {
	decryptedToken, err := u.securityAdapter.Decrypt(request.Token)
	if err != nil {
		u.logs.CustomError("failed to decrypt : %+v", err)
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid email or token")
	}

	emailVerification, err := u.emailVerificationRepo.FindByEmailAndToken(ctx, request.Email, decryptedToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid email or token")
		}
		return helper.WrapInternalServerError(u.logs, "failed to update user email verified_at", err)
	}

	if time.Since(*emailVerification.UpdatedAt) > 15*time.Minute {
		return helper.NewUseCaseError(errorcode.ErrUnauthorized, "Verification token expired")
	}

	now := time.Now()
	user := &entity.User{
		Email: sql.NullString{
			Valid:  true,
			String: request.Email,
		},
		EmailVerifiedAt: &now,
		UpdatedAt:       &now,
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	user, err = u.userRepository.UpdateEmailVerifiedAt(ctx, tx, user)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update user email verified_at", err)
	}

	if err := u.emailVerificationRepo.Delete(ctx, tx, emailVerification); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to delete email verification", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	userProfile, err := u.userProfileRepository.FindByUserId(ctx, user.Id)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to get user profile", err)
	}

	go func() {
		u.createChatRoom(context.Background(), user, userProfile)
		_, _, err := u.createCreatorAndWallet(ctx, user.Id)
		if err != nil {
			u.logs.Error(fmt.Sprintf("Failed to create creator or wallet with error : %v", err))
		}
	}()

	return nil
}

func (u *authUseCase) RequestResetPassword(ctx context.Context, email string) error {
	_, err := u.userRepository.FindByEmailNotGoogle(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid email")
		}
		return helper.WrapInternalServerError(u.logs, "failed to find email", err)
	}

	countByEmailTotal, err := u.resetPasswordRepo.CountByEmail(ctx, email)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed count reset password by email", err)
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	token := uuid.NewString()
	now := time.Now()
	resetPassword := &entity.ResetPassword{
		Email:     email,
		Token:     token,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	if countByEmailTotal > 0 {
		_, err := u.resetPasswordRepo.Update(ctx, tx, resetPassword)
		if err != nil {
			return helper.WrapInternalServerError(u.logs, "failed to update reset password", err)
		}
	} else {
		_, err := u.resetPasswordRepo.Insert(ctx, tx, resetPassword)
		if err != nil {
			return helper.WrapInternalServerError(u.logs, "failed to insert reset password", err)
		}
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	encryptedToken, err := u.securityAdapter.Encrypt(token)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to encrypt reset password token", err)
	}

	if err := u.emailAdapter.SendEmail(email, encryptedToken, "reset password"); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to send reset password email", err)
	}

	return nil
}

func (u *authUseCase) ValidateResetPassword(ctx context.Context, request *model.ValidateResetTokenRequest) (bool, error) {
	decryptedToken, err := u.securityAdapter.Decrypt(request.Token)
	if err != nil {
		return false, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid reset password token")
	}

	_, err = u.resetPasswordRepo.FindByEmailAndToken(ctx, request.Email, decryptedToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "invalid email")
		}
		return false, helper.WrapInternalServerError(u.logs, "failed to find by email and token", err)
	}

	return true, nil
}

func (u *authUseCase) ResetPassword(ctx context.Context, request *model.ResetPasswordUserRequest) error {
	decryptedToken, err := u.securityAdapter.Decrypt(request.Token)
	if err != nil {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid token")
	}

	resetPassword, err := u.resetPasswordRepo.FindByEmailAndToken(ctx, request.Email, decryptedToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid email or token")
		}
		return helper.WrapInternalServerError(u.logs, "failed to find reset password by email and token", err)
	}

	if time.Since(*resetPassword.UpdatedAt) > 15*time.Minute {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Reset password token expired")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	now := time.Now()
	user := &entity.User{
		Email: sql.NullString{
			Valid:  true,
			String: request.Email,
		},
		Password: sql.NullString{
			Valid:  true,
			String: string(hashedPassword),
		},
		UpdatedAt: &now,
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	_, err = u.userRepository.UpdatePassword(ctx, tx, user)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update reset password", err)
	}

	if err := u.resetPasswordRepo.Delete(ctx, tx, resetPassword); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to delete reset password", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}

func (u *authUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, *model.TokenResponse, error) {
	user, err := u.userRepository.FindByMultipleParam(ctx, request.MultipleParam)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, helper.NewUseCaseError(errorcode.ErrValidationFailed, "invalid email or password")
		}
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to send reset password email", err)
	}

	if user.HasVerifiedPhoneNumber() {
		log.Printf("have phone numbers")
	}

	if user.HasEmail() && strings.EqualFold(user.Email.String, request.MultipleParam) {
		if !user.HasVerifiedEmail() {
			return nil, nil, helper.NewUseCaseError(errorcode.ErrValidationFailed, "Email must be verified")
		}
	} else if user.HasPhoneNumber() && user.PhoneNumber.String == request.MultipleParam {
		if !user.HasVerifiedPhoneNumber() {
			return nil, nil, helper.NewUseCaseError(errorcode.ErrValidationFailed, "Phone number must be verified")
		}
	} else {
		return nil, nil, helper.NewUseCaseError(errorcode.ErrValidationFailed, "Invalid login identifier")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(request.Password)); err != nil {
		return nil, nil, helper.NewUseCaseError(errorcode.ErrValidationFailed, "Invalid password")
	}

	var wg sync.WaitGroup
	userProfile, creator, wallet, err := u.getUseProfileWalletAndRepo(ctx, &wg, user.Id)
	if err != nil {
		return nil, nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Invalid email")
	}

	auth := &entity.Auth{
		Id:          user.Id,
		Username:    user.Username,
		Email:       user.Email.String,
		PhoneNumber: user.PhoneNumber.String,
		Similarity:  userProfile.Similarity,
		CreatorId:   creator.Id,
		WalletId:    wallet.Id,
	}

	token, err := u.generateToken(ctx, auth)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to generate token", err)
	}

	return converter.UserToResponse(user), token, nil
}

func (u *authUseCase) generateToken(ctx context.Context, auth *entity.Auth) (*model.TokenResponse, error) {
	accessTokenDetail, err := u.jwtAdapter.GenerateAccessToken(auth.Id)
	if err != nil {
		return nil, fmt.Errorf("generate access token : %+v", err)
	}

	refreshTokenDetail, err := u.jwtAdapter.GenerateRefreshToken(auth.Id)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token : %+v", err)
	}

	jsonValue, err := sonic.ConfigFastest.Marshal(auth)
	if err != nil {
		return nil, fmt.Errorf("marshal user : %+v", err)
	}

	if err := u.cacheAdapter.Set(ctx, refreshTokenDetail.Token, auth.Id, time.Until(refreshTokenDetail.ExpiresAt)); err != nil {
		return nil, fmt.Errorf("save refresh token into cache : %+v", err)
	}

	//TODO -- SYNC with update ()
	//set user persistence data in cache (redis) for better verify auth flow (should be synced with update profile usecase)
	if err := u.cacheAdapter.Set(ctx, auth.Id, jsonValue, time.Until(refreshTokenDetail.ExpiresAt)); err != nil {
		return nil, fmt.Errorf("save user body into cache : %+v", err)
	}

	token := &model.TokenResponse{
		AccessToken:  accessTokenDetail.Token,
		RefreshToken: refreshTokenDetail.Token,
	}

	return token, nil
}

func (u *authUseCase) CreateDeviceToken(ctx context.Context, request *model.DeviceRequest) error {
	now := time.Now()
	userDevice := &entity.UserDevice{
		Id:        ulid.Make().String(),
		UserId:    request.UserId,
		Token:     request.DeviceToken,
		Platform:  request.Platform,
		CreatedAt: &now,
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	_, err = u.userDeviceRepository.Create(ctx, tx, userDevice)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to create user device", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}
	return nil
}

func (u *authUseCase) Current(ctx context.Context, email string) (*model.UserResponse, error) {
	user, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.NewUseCaseError(errorcode.ErrUserNotFound, "Invalid token")
		}
		return nil, helper.WrapInternalServerError(u.logs, "failed to find by email not google", err)
	}

	return converter.UserToResponse(user), nil
}

func (u *authUseCase) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.AuthResponse, error) {
	accessTokenDetail, err := u.jwtAdapter.VerifyAccessToken(request.Token)
	if err != nil {
		return nil, helper.NewUseCaseError(errorcode.ErrUnauthorized, "Invalid access token")
	}

	userId, err := u.cacheAdapter.Get(ctx, request.Token)
	if userId != "" {
		return nil, helper.NewUseCaseError(errorcode.ErrUnauthorized, "User has already signed out")
	}

	cachedUserStr, err := u.cacheAdapter.Get(ctx, accessTokenDetail.UserId)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, helper.WrapInternalServerError(u.logs, "failed to get cached user", err)
	}

	auth := new(entity.Auth)
	//If redis stale, get from db
	if errors.Is(err, redis.Nil) {
		user, err := u.userRepository.FindById(ctx, accessTokenDetail.UserId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, helper.NewUseCaseError(errorcode.ErrUnauthorized, "invalid refresh token")
			}
			return nil, helper.WrapInternalServerError(u.logs, "failed to find user by id", err)
		}

		var wg sync.WaitGroup
		userProfile, creator, wallet, err := u.getUseProfileWalletAndRepo(ctx, &wg, user.Id)
		if err != nil {
			return nil, helper.NewUseCaseError(errorcode.ErrUnauthorized, "invalid refresh token")
		}

		auth = &entity.Auth{
			Id:          user.Id,
			Username:    user.Username,
			Email:       user.Email.String,
			PhoneNumber: user.PhoneNumber.String,
			Similarity:  userProfile.Similarity,
			CreatorId:   creator.Id,
			WalletId:    wallet.Id,
		}

		jsonValue, err := sonic.ConfigFastest.Marshal(auth)
		if err != nil {
			return nil, fmt.Errorf("marshal user : %+v", err)
		}

		if err := u.cacheAdapter.Set(ctx, auth.Id, jsonValue, time.Until(time.Now().Add(time.Hour*24*1))); err != nil {
			return nil, fmt.Errorf("save user body into cache : %+v", err)
		}

	} else {
		if err := sonic.ConfigFastest.Unmarshal([]byte(cachedUserStr), &auth); err != nil {
			return nil, helper.WrapInternalServerError(u.logs, "failed to unmarshal user body from cached", err)
		}
	}

	authResponse := &model.AuthResponse{
		UserId:      auth.Id,
		Username:    auth.Username,
		Email:       auth.Email,
		PhoneNumber: auth.PhoneNumber,
		Similarity:  auth.Similarity,
		CreatorId:   auth.CreatorId,
		WalletId:    auth.WalletId,
		Token:       request.Token,
		ExpiresAt:   accessTokenDetail.ExpiresAt,
	}

	return authResponse, nil
}

func (u *authUseCase) Logout(ctx context.Context, request *model.LogoutUserRequest) (bool, error) {
	if err := u.cacheAdapter.Set(ctx, request.AccessToken, "revoked", time.Until(request.ExpiresAt)); err != nil {
		return false, helper.WrapInternalServerError(u.logs, "failed to save access token to cache for logout : ", err)
	}

	if err := u.cacheAdapter.Del(ctx, request.RefreshToken); err != nil {
		return false, helper.WrapInternalServerError(u.logs, "failed to delete refresh token from cache for logout : ", err)
	}

	return true, nil
}

func (u *authUseCase) AccessTokenRequest(ctx context.Context, refreshToken string) (*model.UserResponse, *model.TokenResponse, error) {
	refreshTokenDetail, err := u.jwtAdapter.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, nil, helper.NewUseCaseError(errorcode.ErrUnauthorized, "Invalid refresh token")
	}

	userId, err := u.cacheAdapter.Get(ctx, refreshToken)
	if userId == "" {
		return nil, nil, helper.NewUseCaseError(errorcode.ErrUnauthorized, "Invalid refresh token")
	}

	user, err := u.userRepository.FindById(ctx, refreshTokenDetail.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, helper.NewUseCaseError(errorcode.ErrUnauthorized, "invalid refresh token")
		}
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to find user by id", err)
	}

	accessTokenDetail, err := u.jwtAdapter.GenerateAccessToken(user.Id)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to generate access token", err)
	}

	tokenResponse := &model.TokenResponse{
		AccessToken: accessTokenDetail.Token,
	}

	//TODO time duration refresh token to be renewed
	if time.Until(refreshTokenDetail.ExpiresAt) < time.Hour*24 {
		newRefreshTokenDetail, err := u.jwtAdapter.GenerateRefreshToken(user.Id)
		if err != nil {
			return nil, nil, helper.WrapInternalServerError(u.logs, "failed to generate refresh token", err)
		}

		if err := u.cacheAdapter.Set(ctx, newRefreshTokenDetail.Token, user.Id, time.Until(refreshTokenDetail.ExpiresAt)); err != nil {
			return nil, nil, helper.WrapInternalServerError(u.logs, "failed to save refresh token to cache", err)
		}

		tokenResponse.RefreshToken = newRefreshTokenDetail.Token
	} else {
		tokenResponse.RefreshToken = refreshTokenDetail.Token
	}

	return converter.UserToResponse(user), tokenResponse, nil
}
