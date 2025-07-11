package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"firebase.google.com/go/v4/messaging"

	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/repository"

	"github.com/redis/go-redis/v9"
)

type NotificationUseCase interface {
	ProcessAndSendBulkNotificationsV2(ctx context.Context, userCountMap map[string]int32) error
	ProcessAndSendSingleFacecamNotifications(ctx context.Context, userID string, countPhotos int) error
	ProcessAndSendSingleNotifications(ctx context.Context, userIDs []string) error
}

type notificationUseCase struct {
	db                    repository.BeginTx
	redisClient           *redis.Client
	userDeviceRepository  repository.UserDeviceRepository
	cloudMessagingAdapter adapter.CloudMessagingAdapter

	logs logger.Log
}

func NewNotificationUseCase(db repository.BeginTx, redisClient *redis.Client, userDeviceRepository repository.UserDeviceRepository,
	cloudMessagingAdapter adapter.CloudMessagingAdapter, logs logger.Log) NotificationUseCase {
	return &notificationUseCase{
		db:                    db,
		redisClient:           redisClient,
		userDeviceRepository:  userDeviceRepository,
		cloudMessagingAdapter: cloudMessagingAdapter,
		logs:                  logs,
	}
}

func (u *notificationUseCase) ProcessAndSendSingleFacecamNotifications(ctx context.Context, userID string, countPhotos int) error {
	log.Println("[USER][NOTIFICATION USECASE]Process and send single facecam notification with userID and count", userID, countPhotos)

	userAuthentications, err := u.fetchFCMTokens(ctx, []string{userID})
	if err != nil {
		log.Println("Error fetching tokens:", err)
		return err
	}
	u.sendFCMWorkerPool(ctx, false, true, userAuthentications, nil, countPhotos, 1)

	return nil
}

func (u *notificationUseCase) ProcessAndSendSingleNotifications(ctx context.Context, userIDs []string) error {
	lenData := len(userIDs)
	log.Println("[USER][NOTIFICATION USECASE]Process and send single notification with len", lenData)

	const batchSize = 5000

	var outerError error
	for i := 0; i < lenData; i += batchSize {
		end := min(i+batchSize, lenData)
		batch := userIDs[i:end]

		userAuthentications, err := u.fetchFCMTokens(ctx, batch)
		if err != nil {
			outerError = err
			log.Println("Error fetching tokens:", err)
			continue
		}
		u.sendFCMWorkerPool(ctx, false, false, userAuthentications, nil, 0, 10)
	}

	if outerError != nil {
		return outerError
	}
	return nil
}

func (u *notificationUseCase) alertAdminFCMAuthIssue(userID string, err *helper.ErrorFCM) {
	u.logs.Log(fmt.Sprintf("[ALERT] Admin needs to investigate auth issue for userID=%s: %s (%s)", userID, err.Code, err.Details))
	// Bisa kirim ke Sentry, Email, atau channel Discord internal kamu
}

func (u *notificationUseCase) getTokenFromRedis(ctx context.Context, userID string) ([]*entity.UserDevice, error) {
	key := fmt.Sprintf("fcm_tokens:%s", userID)
	tokens, err := u.redisClient.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis error for userID %s: %w", userID, err)
	}
	if len(tokens) == 0 {
		return nil, nil
	}

	var devices []*entity.UserDevice
	for _, token := range tokens {
		devices = append(devices, &entity.UserDevice{
			UserId: userID,
			Token:  token,
		})
	}
	return devices, nil
}

// ISSUE Redis Lookup N+1 Problem (Satu per User)
func (u *notificationUseCase) fetchFCMTokens(ctx context.Context, userIDs []string) (*[]*entity.UserDevice, error) {
	u.logs.CustomLog("fetchCMTokens called with userIDS", userIDs)
	if len(userIDs) == 0 {
		return nil, nil
	}

	var (
		missingUserIDs []string
		finalResult    []*entity.UserDevice
	)

	for _, userID := range userIDs {
		devices, err := u.getTokenFromRedis(ctx, userID)
		if err != nil || len(devices) == 0 {
			missingUserIDs = append(missingUserIDs, userID)
			if err != nil {
				u.logs.Log(fmt.Sprintf("Redis error for userID %s: %v", userID, err))
			}
			continue
		}

		finalResult = append(finalResult, devices...)
	}

	// Fallback: ambil dari PostgreSQL untuk user yang tidak ditemukan di Redis
	if len(missingUserIDs) > 0 {
		dbResults, err := u.userDeviceRepository.FetchFCMTokensFromPostgre(ctx, u.db, missingUserIDs)
		if err != nil {
			return nil, fmt.Errorf("postgres fetch error: %w", err)
		}

		// Simpan hasil dari PostgreSQL ke Redis (SADD per user)
		pipe := u.redisClient.Pipeline()
		for _, ua := range *dbResults {
			key := fmt.Sprintf("fcm_tokens:%s", ua.UserId)
			pipe.SAdd(ctx, key, ua.Token)
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			u.logs.Log(fmt.Sprintf("Redis cache update error: %v", err))
		}

		finalResult = append(finalResult, *dbResults...)
	}

	u.logs.Log(fmt.Sprintf("[FETCHFCM][TOTAL_TOKENS] %d", len(finalResult)))
	return &finalResult, nil
}

/* Process and Send Bulk Notification Logic
1. Count user by leveraging concurent computation
2. The total of usersIDS will be divided base on each batch size, then fetch every token
3. Fetch token first done by taking a SMembers Redis (Set Data Structure in redis)
4. For every user id if there are no token, save to missing UserIDs
5. Then 2nd step of fetch token is get missing UserIDs token from postgre (fallback to db)
6. Call Send FCM Worker pool then dispatchWorker using gorotoutine
7. Map every user device token into group by user id then send to channel
8. For every worker, block/wait the usertoken that come from channel
10. Get tokens and count from user device map then send to sendFCMWithRetry
11. Check every cond, if there is business logic error, remove the token
12. If error is external (such as fcm down, ratte limit or max quota) retry is conducted
13. Done
*/

func (u *notificationUseCase) ProcessAndSendBulkNotificationsV2(ctx context.Context, userCountMap map[string]int32) error {
	lenCountMap := len(userCountMap)
	userIDs := make([]string, 0, lenCountMap)
	for idx := range userCountMap {
		userIDs = append(userIDs, idx)
	}

	const batchSize = 5000

	var outerError error
	for i := 0; i < lenCountMap; i += batchSize {
		end := min(i+batchSize, lenCountMap)
		batch := userIDs[i:end]

		userAuthentications, err := u.fetchFCMTokens(ctx, batch)
		if err != nil {
			outerError = err
			log.Println("Error fetching tokens:", err)
			continue
		}
		u.sendFCMWorkerPool(ctx, true, false, userAuthentications, userCountMap, 0, 10)
	}

	if outerError != nil {
		return outerError
	}
	return nil
}

func (u *notificationUseCase) removeUserToken(ctx context.Context, userID, token string) error {
	u.logs.Log(fmt.Sprintf("remove user token called with token : %s and user %s", token, userID))
	if err := repository.BeginTransaction(ctx, u.logs, u.db, func(tx repository.TransactionTx) error {
		if err := u.userDeviceRepository.DeleteByUserIdAndToken(ctx, tx, userID, token); err != nil {
			return err
		}
		u.logs.Log(fmt.Sprintf("Removed FCM token from DB: user=%s, token=%s", userID, token))

		setKey := fmt.Sprintf("fcm_tokens:%s", userID)
		if err := u.redisClient.SRem(ctx, setKey, token).Err(); err != nil {
			u.logs.Log(fmt.Sprintf("Redis SREM error: user=%s, token=%s: %v", userID, token, err))
		} else {
			u.logs.Log(fmt.Sprintf("Removed FCM token from Redis Set: user=%s, token=%s", userID, token))
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (u *notificationUseCase) sendMulticast(ctx context.Context, userID string, tokens []string, message string) error {
	msg := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: "Foto Mirip Terdeteksi",
			Body:  message,
		},
		Data: map[string]string{
			"type":    "similar_photo",
			"message": message,
		},
	}

	res, err := u.cloudMessagingAdapter.SendEachForMulticast(ctx, msg)
	if err != nil {
		return err
	}

	for i, r := range res.Responses {
		if !r.Success {
			token := tokens[i]
			u.logs.Log(fmt.Sprintf("[MULTICAST][CHECKER] userID=%s, token=%s, raw=%v", userID, token, r.Error))

			fcmErr := helper.ParseFCMError(r.Error)

			if fcmErr.IsInvalidToken() {
				u.logs.Log(fmt.Sprintf("[MULTICAST][INVALID] userID=%s, token=%s, code=%s raw = %s", userID, token, fcmErr.Code, fcmErr.Raw))
				_ = u.removeUserToken(ctx, userID, token)
			} else {
				u.logs.Log(fmt.Sprintf("[MULTICAST][FAIL] userID=%s, token=%s, code=%s, detail=%s", userID, token, fcmErr.Code, fcmErr.Details))
				_ = u.removeUserToken(ctx, userID, token)
			}
		}
	}

	u.logs.Log(fmt.Sprintf("✅ Sent multicast to userID=%s, success=%d, failed=%d", userID, res.SuccessCount, res.FailureCount))
	return nil
}

func (u *notificationUseCase) sendFCMMulticastWithRetry(ctx context.Context, userID string, tokens []string, message string) error {
	const maxRetries = 3
	backoff := 500 * time.Millisecond

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := u.sendMulticast(ctx, userID, tokens, message)
		if err == nil {
			return nil
		}

		fcmErr := helper.ParseFCMError(err)
		if fcmErr.IsRetryable() {
			u.logs.Log(fmt.Sprintf("[MULTICAST][RETRY] userID=%s, retry attempt=%d, err=%v", userID, attempt, err))
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		if fcmErr.IsAuthError() {
			u.logs.Log(fmt.Sprintf("[AUTH ERROR] userID=%s, code=%s → admin action required", userID, fcmErr.Code))
			u.alertAdminFCMAuthIssue(userID, fcmErr)
			return fmt.Errorf("auth error sending FCM to userID=%s: %w", userID, err)
		}

		u.logs.Log(fmt.Sprintf("[MULTICAST][FATAL ERROR] userID=%s, err=%v", userID, err))
		return err
	}

	return fmt.Errorf("failed multicast after retries for userID=%s", userID)
}

func (u *notificationUseCase) sendFCMWorkerPool(ctx context.Context, isBulk bool, isSingleFacecam bool, userAuthentications *[]*entity.UserDevice, countMap map[string]int32, countPhotos int, workerCount int) {
	// 1. Group tokens per userID
	userTokensMap := make(map[string][]string)
	for _, ua := range *userAuthentications {
		userTokensMap[ua.UserId] = append(userTokensMap[ua.UserId], ua.Token)
	}

	// 2. Job channel per userID
	jobChan := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case userID, ok := <-jobChan:
					if !ok {
						return
					}

					tokens := userTokensMap[userID]
					var count int32 = 1
					if isBulk {
						count = countMap[userID]
						if count == 0 {
							continue
						}
					}

					if isSingleFacecam {
						count = int32(countPhotos)
					}

					if len(tokens) == 0 {
						continue
					}

					message := fmt.Sprintf("Terdapat %d foto yang mirip dengan Anda!", count)
					if err := u.sendFCMMulticastWithRetry(ctx, userID, tokens, message); err != nil {
						u.logs.Log(fmt.Sprintf("[Worker %d] ❌ Failed send to userID=%s: %v", workerID, userID, err))
					}
				}
			}
		}(i)
	}

	// 3. Kirim pekerjaan ke worker
	go func() {
		for userID := range userTokensMap {
			select {
			case <-ctx.Done():
				close(jobChan)
				return
			case jobChan <- userID:
			}
		}
		close(jobChan)
	}()

	// 4. Tunggu semua worker selesai
	wg.Wait()
}
