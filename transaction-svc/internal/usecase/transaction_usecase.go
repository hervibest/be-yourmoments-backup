package usecase

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/enum"
	errorcode "github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/enum/error"
	producer "github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/gateway/messaging"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model/converter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model/event"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/repository"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/usecase/contract"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type transactionUseCase struct {
	db           *sqlx.DB
	photoAdapter adapter.PhotoAdapter

	transactionRepository       repository.TransactionRepository
	transactionDetailRepository repository.TransactionDetailRepository
	transactionItemRepository   repository.TransactionItemRepository
	walletRepository            repository.WalletRepository
	transactionWalletRepo       repository.TransactionWalletRepository

	paymentAdapter      adapter.PaymentAdapter
	cacheAdapter        adapter.CacheAdapter
	timeParserHelper    helper.TimeParserHelper
	transactionProducer producer.TransactionProducer
	logs                *logger.Log
}

func NewTransactionUseCase(db *sqlx.DB, photoAdapter adapter.PhotoAdapter, transactionRepository repository.TransactionRepository,
	transactionItemRepository repository.TransactionItemRepository, transactionDetailRepository repository.TransactionDetailRepository,
	walletRepository repository.WalletRepository, transactionWalletRepo repository.TransactionWalletRepository,
	paymentAdapter adapter.PaymentAdapter, cacheAdapter adapter.CacheAdapter, timeParserHelper helper.TimeParserHelper,
	transactionProducer producer.TransactionProducer, logs *logger.Log) contract.TransactionUseCase {

	if db == nil {
		logs.Error("Sqlx.DB is nil")
	}

	return &transactionUseCase{
		db:                          db,
		photoAdapter:                photoAdapter,
		transactionRepository:       transactionRepository,
		transactionItemRepository:   transactionItemRepository,
		transactionDetailRepository: transactionDetailRepository,
		walletRepository:            walletRepository,
		transactionWalletRepo:       transactionWalletRepo,
		paymentAdapter:              paymentAdapter,
		cacheAdapter:                cacheAdapter,
		timeParserHelper:            timeParserHelper,
		transactionProducer:         transactionProducer,
		logs:                        logs,
	}
}

func (u *transactionUseCase) CreateTransaction(ctx context.Context, request *model.CreateTransactionRequest) (*model.CreateTransactionResponse, error) {
	items, total, err := u.photoAdapter.CalculatePhotoPrice(ctx, request.UserId, request.CreatorId, request.PhotoIds)
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
		Id:             uuid.NewString(),
		UserId:         request.UserId,
		InternalStatus: enum.TrxInternalStatusPending,
		Status:         enum.TransactionStatusPending,
		PhotoIds:       photoIDsByte,
		Amount:         total.Price,
		CheckoutAt:     &now,
		CreatedAt:      &now,
		UpdatedAt:      &now,
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
				CreatedAt:           &now,
				UpdatedAt:           &now,
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
		return converter.TransactionToResponse(transaction, ""), nil
	}

	if err := u.updateTransactionToken(ctx, token, transaction.Id); err != nil {
		return nil, err
	}

	if err = u.transactionProducer.ScheduleTransactionTaskExpiration(ctx, transaction.Id); err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to set transaction task expiration", err)
	}

	transaction.SnapToken.String = token

	return converter.TransactionToResponse(transaction, redirectUrl), nil
}

func (u *transactionUseCase) updateTransactionToken(ctx context.Context, token, transactionId string) error {
	u.logs.Log(fmt.Sprintf("Update transaction token :%s with trx id : %s", token, transactionId))
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	now := time.Now()
	transaction := &entity.Transaction{
		Id:             transactionId,
		InternalStatus: enum.TrxInternalStatusTokenReady,
		SnapToken:      sql.NullString{String: token, Valid: true},
		UpdatedAt:      &now,
	}
	u.logs.Log(fmt.Sprintf("Update transaction transaction %v", transaction))

	err = u.transactionRepository.UpdateToken(ctx, tx, transaction)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update snapshot token in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}
	return nil
}

func (u *transactionUseCase) getPaymentToken(ctx context.Context, transaction *entity.Transaction) (string, string, error) {
	snapRequest := &model.PaymentSnapshotRequest{
		OrderID:     transaction.Id,
		GrossAmount: int64(transaction.Amount),
		Email:       "",
	}

	snapResponse, err := u.paymentAdapter.CreateSnapshot(ctx, snapRequest)
	if err == nil {
		return snapResponse.Token, snapResponse.RedirectURL, nil
	}

	// Jalankan retry async jika terjadi error
	go func() {
		const maxRetry = 5
		retry := 0

		for retry < maxRetry && !errors.Is(err, gobreaker.ErrOpenState) {
			time.Sleep(time.Second * time.Duration(2<<retry)) // exponential backoff: 2s, 4s, 8s, ...

			snapResponse, retryErr := u.paymentAdapter.CreateSnapshot(context.TODO(), snapRequest)
			if retryErr == nil {
				u.logs.CustomLog("snap response :", snapResponse)
				u.logs.Log(fmt.Sprintf("[RetrySuccess] Transaction %s berhasil mendapatkan snap token setelah %d retry", transaction.Id, retry))
				if err := u.updateTransactionToken(ctx, snapResponse.Token, transaction.Id); err != nil {
					u.logs.Error(fmt.Sprintf("[UpdateFailed] Update transaction token failed with reason : %v", err))
				}
				return
			}

			if errors.Is(retryErr, gobreaker.ErrOpenState) {
				u.logs.Error(fmt.Sprintf("[BreakerOpen] Retry stopped at %d, circuit breaker open", retry))
				break
			}

			err = retryErr
			retry++
		}

		// Jika tetap gagal karena breaker open, listen redis stream
		u.logs.Log(fmt.Sprintf("[WaitRecovery] Transaction %s menunggu circuit breaker pulih via redis", transaction.Id))
		u.subscribeMidtransRecoveryAndRetry(u.paymentAdapter, snapRequest, transaction.Id)
	}()

	return "", "", fmt.Errorf("midtrans create snapshot error: %w", err)
}

func (u *transactionUseCase) subscribeMidtransRecoveryAndRetry(adapter adapter.PaymentAdapter, req *model.PaymentSnapshotRequest, txId string) {
	for {
		streams, err := u.cacheAdapter.XRead(context.Background(), &redis.XReadArgs{
			Streams: []string{"midtrans:recovery", "0"},
			Block:   0, // block selamanya sampai ada signal
			Count:   1,
		})

		if err != nil {
			log.Printf("[StreamError] Redis read error: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}

		for _, msg := range streams[0].Messages {
			log.Printf("[RecoverySignal] Received for tx %s: %+v", txId, msg.Values)
			ctx := context.Background()
			resp, err := adapter.CreateSnapshot(ctx, req)
			if err != nil {
				log.Printf("[RecoveryFail] Retry gagal lagi setelah sinyal pulih: %v", err)
				continue
			}

			log.Printf("[RecoverySuccess] Retry berhasil: tx %s → token %s", txId, resp.Token)
			if err := u.updateTransactionToken(ctx, resp.Token, txId); err != nil {
				u.logs.Error(fmt.Sprintf("[UpdateFailed] Update transaction token failed with reason : %v", err))
			}
			return
		}
	}
}

func (u *transactionUseCase) CheckPaymentSignature(signatureKey, transcationId, statusCode, grossAmount string) (bool, string) {
	u.logs.Log(fmt.Sprintf("Ini adalah transactionId : %s", transcationId))
	u.logs.Log(fmt.Sprintf("Ini adalah status code : %s", statusCode))
	u.logs.Log(fmt.Sprintf("Ini adalah gross amount : %s", grossAmount))
	u.logs.Log(fmt.Sprintf("Ini adalah signature key : %s", signatureKey))

	signatureToCompare := transcationId + statusCode + grossAmount + u.paymentAdapter.GetPaymentServerKey()
	hash := sha512.New()
	hash.Write([]byte(signatureToCompare))
	hashedSignature := hex.EncodeToString(hash.Sum(nil))
	return hashedSignature == signatureKey, hashedSignature
}

func (u *transactionUseCase) CheckAndUpdateTransaction(ctx context.Context, request *model.CheckAndUpdateTransactionRequest) error {
	if u.db == nil {
		fmt.Print("DB NIL")
		u.logs.Error("Sqlx.DB is nil")
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	transaction, err := u.transactionRepository.FindById(ctx, tx, request.OrderID, true)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update transaction callback in database", err)
	}

	requestIsValid, hashedSignature := u.CheckPaymentSignature(request.SignatureKey, transaction.Id, request.StatusCode, request.GrossAmount)
	if !requestIsValid {
		u.logs.Log(fmt.Sprintf("Invalid signature: expected=%s got=%s", hashedSignature, request.SignatureKey))
		return helper.NewUseCaseError(errorcode.ErrForbidden, "Invalid signature key")
	}

	//Only check if internal transaction status is "PENDING", "TOKEN_READY", "STATUS_EXPIRED"
	if transaction.InternalStatus != enum.TrxInternalStatusPending &&
		transaction.InternalStatus != enum.TrxInternalStatusTokenReady &&
		transaction.InternalStatus != enum.TrxInternalStatusExpired {
		u.logs.CustomLog("transaction internal status already final. Ignoring duplicate settlement :.", transaction.Id)
		return nil
	}

	now := time.Now()
	updateTransaction := &entity.Transaction{
		Id:        transaction.Id,
		UpdatedAt: &now,
	}

	var transactionStatus enum.TransactionStatus
	var transactionInternalStatus enum.TrxInternalStatus
	var settlementTimePtr *time.Time

	//First internal expired status checking
	if transaction.InternalStatus == enum.TrxInternalStatusExpired {
		// IF External Or Midtrans Payment Settled Check the Settlement Time
		if request.MidtransTransactionStatus == string(enum.PaymentStatusSettlement) {
			settlementTime, err := u.timeParserHelper.TimeParseInDefaultLocation(request.SettlementTime)
			if err != nil {
				return err
			}

			settlementTimePtr = &settlementTime
			graceDeadline := transaction.UpdatedAt.Add(5 * time.Minute)

			//Checking settlement time from external midtrans with internal expired time
			if settlementTime.After(graceDeadline) {
				transactionStatus = enum.TransactionStatusExpired //User transaction is valid or settled
				transactionInternalStatus = enum.TrxInternalStatusLateSettlement
			} else {
				//User transaction is valid or settled even the user come late but from system not the settled time
				transactionStatus = enum.TransactionStatusSuccess
				transactionInternalStatus = enum.TrxInternalStatusExpiredCheckedValid
				updateTransaction.PaymentAt = &now
			}
			updateTransaction.ExternalSettlementAt = settlementTimePtr
			updateTransaction.SnapToken = sql.NullString{String: "", Valid: true}
		} else {
			transactionInternalStatus = enum.TrxInternalStatusExpiredCheckedInvalid //User doesnt settled even internal status has expired
			transactionStatus = enum.TransactionStatusExpired
		}

	} else {
		switch request.MidtransTransactionStatus {
		case string(enum.PaymentStatusCapture), string(enum.PaymentStatusSettlement):
			settlementTime, err := u.timeParserHelper.TimeParseInDefaultLocation(request.SettlementTime)
			if err != nil {
				return err
			}
			settlementTimePtr = &settlementTime
			transactionStatus = enum.TransactionStatusSuccess
			transactionInternalStatus = enum.TrxInternalStatusSettled
			updateTransaction.PaymentAt = &now
			updateTransaction.SnapToken = sql.NullString{String: "", Valid: true}
			updateTransaction.ExternalSettlementAt = settlementTimePtr

		case string(enum.PaymentStatusPending):
			transactionStatus = enum.TransactionStatusPending
			transactionInternalStatus = enum.TrxInternalStatusPending

		case string(enum.PaymentStatusExpire):
			transactionStatus = enum.TransactionStatusExpired
			transactionInternalStatus = enum.TrxInternalStatusExpiredCheckedInvalid
			updateTransaction.SnapToken = sql.NullString{String: "", Valid: true}

		case string(enum.PaymentStatusFailure), string(enum.PaymentStatusDeny):
			transactionStatus = enum.TransactionStatusFailed
			transactionInternalStatus = enum.TrxInternalStatusFailed
			updateTransaction.SnapToken = sql.NullString{String: "", Valid: true}

		case string(enum.PaymentStatusCancel):
			transactionStatus = enum.TransactionStatusCancelled
			transactionInternalStatus = enum.TrxInternalStatusCancelledBySystem
			updateTransaction.SnapToken = sql.NullString{String: "", Valid: true}
		}
	}

	updateTransaction.Status = transactionStatus
	updateTransaction.InternalStatus = transactionInternalStatus
	updateTransaction.ExternalStatus = sql.NullString{
		Valid:  true,
		String: request.MidtransTransactionStatus,
	}

	externalCallbackResponse := json.RawMessage(request.Body)
	updateTransaction.ExternalCallbackResponse = &externalCallbackResponse

	if err := u.transactionRepository.UpdateCallback(ctx, tx, updateTransaction); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to update transaction callback in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	var photoIds []string

	if err := sonic.ConfigFastest.Unmarshal([]byte(transaction.PhotoIds), &photoIds); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to unmarshal photo ids", err)
	}

	if transactionStatus == enum.TransactionStatusCancelled || transactionStatus == enum.TransactionStatusExpired ||
		transactionStatus == enum.TransactionStatusFailed {
		// if err := u.photoAdapter.CancelPhotos(ctx, transaction.UserId, photoIds); err != nil {
		// 	return err
		// }

		event := &event.CancelPhotosEvent{
			UserId:   transaction.UserId,
			PhotoIds: photoIds,
		}

		if err := u.transactionProducer.ProduceTransactionCanceledEvent(ctx, event); err != nil {
			u.logs.Error(fmt.Sprintf("failed to publish cancel photos event: %v", err))
			return nil
		}

		u.logs.Log(fmt.Sprintf("successfully cancel photos from photo service with transaction id %s and buyer %s :.", transaction.Id, transaction.UserId))
		return nil
	}

	if transactionStatus != enum.TransactionStatusSuccess {
		u.logs.CustomLog("transaction not yet success. With external transaction status :.", transaction.ExternalStatus)
		return nil
	}

	// //Call photo service to update photo owner
	// if err := u.photoAdapter.OwnerOwnPhotos(ctx, transaction.UserId, photoIds); err != nil {
	// 	return err
	// }

	event := &event.OwnerOwnPhotosEvent{
		UserId:   transaction.UserId,
		PhotoIds: photoIds,
	}

	if err := u.transactionProducer.ProduceTransactionSettledEvent(ctx, event); err != nil {
		u.logs.Error(fmt.Sprintf("failed to publish owner own settled photos event: %v", err))
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
		return fmt.Errorf("find many transaction detail by trx id error : %+v", err)
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
		return fmt.Errorf("find wallets by creator ids error : %+v", err)
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
		return fmt.Errorf("failed to bulk insert transaction wallet :  %+v", err)
	}

	// 6. Update wallet balance
	for walletID, addAmount := range walletUpdateMap {
		if err := u.walletRepository.AddBalance(ctx, tx, walletID, int64(addAmount)); err != nil {
			return fmt.Errorf("failed to add balance to wallet :  %+v", err)
		}
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}

func (u *transactionUseCase) UserGetWithDetail(ctx context.Context, request *model.GetTransactionWithDetail) (*model.TransactionWithDetail, error) {
	transactionWithDetails, err := u.transactionRepository.UserFindWithDetailById(ctx, u.db, request.TransactionId, request.UserID)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to get user transaction with detail by id", err)
	}

	if len(*transactionWithDetails) == 0 {
		return nil, helper.NewUseCaseError(errorcode.ErrResourceNotFound, "invalid transaction id")
	}

	var photoIds []string
	if err := sonic.ConfigFastest.Unmarshal([]byte((*transactionWithDetails)[0].PhotoIds), &photoIds); err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to unmarshal photo ids", err)
	}

	pbPhotoWithDetails, err := u.photoAdapter.GetPhotoWithDetails(ctx, photoIds, request.UserID)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to get photo with details from photo service using grpc", err)
	}

	return converter.TransactionAndPhotoToSingleResponse(*transactionWithDetails, *pbPhotoWithDetails), nil
}

func (u *transactionUseCase) GetAllUserTransaction(ctx context.Context, request *model.GetAllUsertTransaction) (*[]*model.UserTransaction, *model.PageMetadata, error) {
	userTransactions, pageMetadata, err := u.transactionRepository.UserFindAll(ctx, u.db, request.Page, request.Size, request.UserId, request.Order)
	if err != nil {
		return nil, nil, helper.WrapInternalServerError(u.logs, "failed to find all creator review", err)
	}

	return converter.UserTransactionToResponse(userTransactions), pageMetadata, nil
}

func (u *transactionUseCase) CreateTransactionV2(ctx context.Context, request *model.CreateTransactionV2Request) (*model.CreateTransactionResponse, error) {
	var outerErr error

	items, total, errAdapter := u.photoAdapter.CalculatePhotoPriceV2(ctx, request.UserId, request.CreatorId, request)
	if errAdapter != nil {
		return nil, errAdapter
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

	photoIDs := make([]string, 0, len(*items))

	for _, item := range *items {
		photoIDs = append(photoIDs, item.PhotoId)
	}

	photoIDsByte, err := sonic.ConfigFastest.Marshal(photoIDs)
	if err != nil {
		outerErr = err
		return nil, helper.WrapInternalServerError(u.logs, "failed to marshal photo ids", err)
	}

	transaction := &entity.Transaction{
		Id:             uuid.NewString(),
		UserId:         request.UserId,
		InternalStatus: enum.TrxInternalStatusPending,
		Status:         enum.TransactionStatusPending,
		PhotoIds:       photoIDsByte,
		Amount:         total.Price,
		CheckoutAt:     &now,
		CreatedAt:      &now,
		UpdatedAt:      &now,
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
				CreatedAt:           &now,
				UpdatedAt:           &now,
			}
			allItems = append(allItems, &item)
		}
	}

	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		outerErr = err
		return nil, err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	transaction, err = u.transactionRepository.Create(ctx, tx, transaction)
	if err != nil {
		outerErr = err
		return nil, helper.WrapInternalServerError(u.logs, "failed to create transaction in database", err)
	}

	_, err = u.transactionDetailRepository.Create(ctx, tx, transaction.Id, details)
	if err != nil {
		outerErr = err
		return nil, helper.WrapInternalServerError(u.logs, "failed to create transaction details in database", err)
	}

	_, err = u.transactionItemRepository.Create(ctx, tx, allItems)
	if err != nil {
		outerErr = err
		return nil, helper.WrapInternalServerError(u.logs, "failed to create transaction items in database", err)
	}

	if err := repository.Commit(tx, u.logs); err != nil {
		outerErr = err
		return nil, err
	}

	defer func() {
		if outerErr != nil {
			event := &event.CancelPhotosEvent{
				UserId:   request.UserId,
				PhotoIds: photoIDs,
			}
			if err := u.transactionProducer.ProduceTransactionCanceledEvent(ctx, event); err != nil {
				u.logs.Error(fmt.Sprintf("failed to publish cancel photos event: %v", err))
			}
		}
	}()

	token, redirectUrl, err := u.getPaymentToken(ctx, transaction)
	if err != nil {
		return converter.TransactionToResponse(transaction, ""), nil
	}

	if err := u.updateTransactionToken(ctx, token, transaction.Id); err != nil {
		return nil, err
	}

	if err = u.transactionProducer.ScheduleTransactionTaskExpiration(ctx, transaction.Id); err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to set transaction task expiration", err)
	}

	transaction.SnapToken.String = token

	return converter.TransactionToResponse(transaction, redirectUrl), nil
}
