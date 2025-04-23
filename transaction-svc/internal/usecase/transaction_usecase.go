package usecase

import (
	"be-yourmoments/transaction-svc/internal/adapter"
	"be-yourmoments/transaction-svc/internal/entity"
	"be-yourmoments/transaction-svc/internal/enum"
	errorcode "be-yourmoments/transaction-svc/internal/enum/error"
	"be-yourmoments/transaction-svc/internal/helper"
	"be-yourmoments/transaction-svc/internal/helper/logger"
	"be-yourmoments/transaction-svc/internal/model"
	"be-yourmoments/transaction-svc/internal/model/converter"
	"be-yourmoments/transaction-svc/internal/repository"
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, request *model.CreateTransactionRequest) (*model.CreateTransactionResponse, error)
	UpdateTransactionWebhook(ctx context.Context, request *model.UpdateTransactionWebhookRequest) error
}

type transactionUseCase struct {
	db           *sqlx.DB
	photoAdapter adapter.PhotoAdapter

	transactionRepository       repository.TransactionRepository
	transactionDetailRepository repository.TransactionDetailRepository
	transactionItemRepository   repository.TransactionItemRepository

	paymentAdapter adapter.PaymentAdapter
	logs           *logger.Log

	walletRepository      repository.WalletRepository
	transactionWalletRepo repository.TransactionWalletRepository
}

func NewTransactionUseCase(db *sqlx.DB, photoAdapter adapter.PhotoAdapter, transactionRepository repository.TransactionRepository,
	transactionItemRepository repository.TransactionItemRepository, transactionDetailRepository repository.TransactionDetailRepository,
	paymentAdapter adapter.PaymentAdapter, logs *logger.Log, walletRepository repository.WalletRepository,
	transactionWalletRepo repository.TransactionWalletRepository) TransactionUseCase {
	return &transactionUseCase{
		db:                          db,
		photoAdapter:                photoAdapter,
		transactionRepository:       transactionRepository,
		transactionItemRepository:   transactionItemRepository,
		transactionDetailRepository: transactionDetailRepository,
		paymentAdapter:              paymentAdapter,
		logs:                        logs,
		walletRepository:            walletRepository,
		transactionWalletRepo:       transactionWalletRepo,
	}
}

func (u *transactionUseCase) CreateTransaction(ctx context.Context, request *model.CreateTransactionRequest) (*model.CreateTransactionResponse, error) {
	items, total, err := u.photoAdapter.CalculatePhotoPrice(ctx, request.UserId, request.PhotoIds)
	if err != nil {
		return nil, err
	}

	if len(*items) == 0 {
		return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Items is empty")
	}

	if total.Price == 0 {
		return nil, helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Total price is zero")
	}

	// Step 1: Grouping per creator (photographer)
	grouped := map[string][]*model.CheckoutItem{}
	for _, i := range *items {
		grouped[i.CreatorId] = append(grouped[i.CreatorId], i)
	}

	// Step 2: Hitung subtotal per creator dan total amount
	var totalAmount int32 = 0

	details := make([]*entity.TransactionDetail, 0)

	allItems := make([]*entity.TransactionItem, 0)
	now := time.Now()

	photoIDsByte, err := sonic.ConfigFastest.Marshal(request.PhotoIds)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to marshal photo ids", err)
	}

	transaction := &entity.Transaction{
		Id:         uuid.NewString(),
		UserId:     request.UserId,
		Status:     enum.TransactionStatusPending,
		PhotoIds:   photoIDsByte,
		Amount:     total.Price,
		CheckoutAt: &now,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	for creatorID, list := range grouped {
		var subtotal int32 = 0
		for _, it := range list {
			subtotal += it.FinalPrice
		}

		detail := entity.TransactionDetail{
			Id:                ulid.Make().String(),
			TransactionId:     transaction.Id,
			CreatorId:         creatorID,
			SubTotalPrice:     subtotal,
			CreatorDiscountId: list[0].DiscountId,
			CreatedAt:         &now,
			UpdatedAt:         &now,
		}

		details = append(details, &detail)
		totalAmount += subtotal

		for _, it := range list {
			item := entity.TransactionItem{
				Id:                  ulid.Make().String(),
				TransactionDetailId: detail.Id,
				PhotoId:             it.PhotoId,
				Price:               it.Price,
				Discount:            sql.NullInt32{Int32: it.Discount, Valid: it.Discount != 0},
				FinalPrice:          it.FinalPrice,
				CreatedAt:           now,
				UpdatedAt:           now,
			}
			allItems = append(allItems, &item)
		}
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return nil, err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	transaction, err = u.transactionRepository.Create(ctx, tx, transaction)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to create transaction in database", err)
	}

	_, err = u.transactionDetailRepository.Create(ctx, tx, transaction.Id, details)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to create transaction details in database", err)
	}

	_, err = u.transactionItemRepository.Create(ctx, tx, allItems)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to create transaction items in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return nil, err
	}

	token, redirectUrl, err := u.getPaymentToken(ctx, transaction)
	if err != nil {
		return nil, helper.WrapExternalServiceUnavailable(u.logs, "failed to get payment token from midtrans", err)
	}

	tx, err = repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return nil, err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	now = time.Now()
	transaction = &entity.Transaction{
		Id:        transaction.Id,
		SnapToken: sql.NullString{String: token},
		UpdatedAt: &now,
	}

	err = u.transactionRepository.UpdateToken(ctx, tx, transaction)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to update snapshot token in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return nil, err
	}

	return converter.TransactionToResponse(transaction, redirectUrl), nil
}

func (u *transactionUseCase) getPaymentToken(ctx context.Context, transaction *entity.Transaction) (string, string, error) {
	snapRequest := &model.PaymentSnapshotRequest{
		OrderID:     transaction.Id,
		GrossAmount: int64(transaction.Amount),
		Email:       "",
	}

	snapResponse, err := u.paymentAdapter.CreateSnapshot(ctx, snapRequest)
	if err != nil {
		return "", "", fmt.Errorf("midtrans create snapshot error : %w", err)
	}

	return snapResponse.Token, snapResponse.RedirectURL, nil
}

func (u *transactionUseCase) UpdateTransactionWebhook(ctx context.Context, request *model.UpdateTransactionWebhookRequest) error {
	transaction, err := u.transactionRepository.FindById(ctx, u.db, request.OrderID)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update transaction callback in database", err)
	}

	signatureToCompare := transaction.Id + request.StatusCode + request.GrossAmount + u.paymentAdapter.GetPaymentServerKey()

	hash := sha512.New()
	hash.Write([]byte(signatureToCompare))
	hashedSignature := hex.EncodeToString(hash.Sum(nil))

	requestIsValid := hashedSignature == request.SignatureKey
	if !requestIsValid {
		return helper.NewUseCaseError(errorcode.ErrForbidden, "Invalid signature key")
	}

	midtransTransactionStatus := enum.MidtransPaymentStatus(request.TransactionStatus)
	if midtransTransactionStatus != enum.PaymentStatusSettlement {
		return nil
	}

	if transaction.Status == enum.TransactionStatusSuccess {
		u.logs.CustomLog("transaction already success. Ignoring duplicate settlement webhook :.", transaction.Id)
		return nil
	}

	updateTransaction := &entity.Transaction{
		Id: transaction.Id,
	}

	ExternalCallbackResponse := json.RawMessage(request.Body)

	var transactionStatus enum.TransactionStatus

	now := time.Now()
	switch request.TransactionStatus {
	case string(enum.PaymentStatusSettlement):
		transactionStatus = enum.TransactionStatusSuccess
		updateTransaction.PaymentAt = &now
		updateTransaction.UpdatedAt = &now
		updateTransaction.SnapToken = sql.NullString{String: "", Valid: true}
	case string(enum.PaymentStatusPending):
		transactionStatus = enum.TransactionStatusPending
	case string(enum.PaymentStatusExpire):
		transactionStatus = enum.TransactionStatusExpired
	case string(enum.PaymentStatusFailure):
		transactionStatus = enum.TransactionStatusFailed
	}

	updateTransaction.Status = transactionStatus
	updateTransaction.ExternalStatus = sql.NullString{
		Valid:  true,
		String: request.TransactionStatus,
	}
	updateTransaction.ExternalCallbackResponse = &ExternalCallbackResponse

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	if err := u.transactionRepository.UpdateCallback(ctx, tx, updateTransaction); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update transaction callback in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	if transactionStatus != enum.TransactionStatusSuccess {
		u.logs.CustomLog("transaction not yet success. With external transaction status :.", transaction.ExternalStatus)
		return nil
	}

	var photoIds []string

	if err := sonic.ConfigFastest.Unmarshal([]byte(transaction.PhotoIds), &photoIds); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to unmarshal photo ids", err)
	}

	//Call photo service to update photo owner
	if err := u.photoAdapter.OwnerOwnPhotos(ctx, transaction.UserId, photoIds); err != nil {
		return err
	}

	if err := u.distributeTransactionToWallets(ctx, transaction.Id); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to distribute transaction to wallets", err)
	}

	return nil
}

func (u *transactionUseCase) distributeTransactionToWallets(ctx context.Context, transactionID string) error {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	// 1. Get all transaction_details for the transaction_id
	details, err := u.transactionDetailRepository.FindManyByTrxID(ctx, tx, transactionID)
	if err != nil {
		return errors.New(fmt.Sprintf("find many transaction detail by trx id error : %+v", err))
	}

	if len(*details) == 0 {
		return nil // no detail, nothing to do
	}

	// 2. Get all creator_ids involved
	creatorIDMap := map[string]struct{}{}
	for _, d := range *details {
		creatorIDMap[d.CreatorId] = struct{}{}
	}

	// Convert map keys to slice
	creatorIDs := make([]string, 0, len(creatorIDMap))
	for id := range creatorIDMap {
		creatorIDs = append(creatorIDs, id)
	}

	// 3. Get wallets by creator_id
	wallets, err := u.walletRepository.FindByCreatorIDs(ctx, tx, creatorIDs)
	if err != nil {
		return errors.New(fmt.Sprintf("find wallets by creator ids error : %+v", err))
	}

	// Mapping creator_id -> wallet_id
	walletMap := map[string]string{}
	for _, w := range *wallets {
		walletMap[w.CreatorId] = w.Id
	}

	// 4. Prepare bulk insert for transaction_wallets
	txWallets := make([]*entity.TransactionWallet, 0)
	walletUpdateMap := map[string]int{}

	now := time.Now()
	for _, d := range *details {
		walletID := walletMap[d.CreatorId]
		txWallets = append(txWallets, &entity.TransactionWallet{
			Id:                  ulid.Make().String(),
			WalletId:            walletID,
			TransactionDetailId: d.Id,
			Amount:              d.SubTotalPrice,
			CreatedAt:           &now,
			UpdatedAt:           &now,
		})
		walletUpdateMap[walletID] += int(d.SubTotalPrice)
	}

	// 5. Bulk insert transaction_wallets
	if err := u.transactionWalletRepo.BulkInsert(ctx, tx, txWallets); err != nil {
		return errors.New(fmt.Sprintf("failed to bulk insert transaction wallet :  %+v", err))
	}

	// 6. Update wallet balance
	for walletID, addAmount := range walletUpdateMap {
		if err := u.walletRepository.AddBalance(ctx, tx, walletID, int64(addAmount)); err != nil {
			return errors.New(fmt.Sprintf("failed to add balance to wallet :  %+v", err))
		}
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}
