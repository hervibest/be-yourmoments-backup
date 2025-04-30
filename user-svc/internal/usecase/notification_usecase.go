package usecase

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/repository"
	"google.golang.org/api/googleapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"firebase.google.com/go/messaging"
	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type NotificationUseCase interface {
	ProcessAndSendNotifications(ctx context.Context, datas []*photopb.BulkUserSimilarPhoto) error
}

type notificationUseCase struct {
	db                    *sqlx.DB
	redisClient           *redis.Client
	userDeviceRepository  repository.UserDeviceRepository
	cloudMessagingAdapter adapter.CloudMessagingAdapter

	logs *logger.Log
}

func NewNotificationUseCase(db *sqlx.DB, redisClient *redis.Client, userDeviceRepository repository.UserDeviceRepository,
	cloudMessagingAdapter adapter.CloudMessagingAdapter, logs *logger.Log) NotificationUseCase {
	return &notificationUseCase{
		db:                    db,
		redisClient:           redisClient,
		userDeviceRepository:  userDeviceRepository,
		cloudMessagingAdapter: cloudMessagingAdapter,
		logs:                  logs,
	}
}

// --- Fungsi untuk parallel counting users dari photos ---
// pemanfaatan pararel computation untuk handle I/O bottleneck
// proyeksi kalau creator upload 100 foto sekaligus dan each poto bisa sekitar 100 hingga 1000 user
func (u *notificationUseCase) countUsersParallel(datas []*photopb.BulkUserSimilarPhoto) map[string]int {
	countMap := make(map[string]int)
	var mu sync.Mutex

	numWorkers := runtime.NumCPU()
	chunkSize := (len(datas) + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := min(start+chunkSize, len(datas))
		if start >= len(datas) {
			continue
		}

		wg.Add(1)
		go func(part []*photopb.BulkUserSimilarPhoto) {
			defer wg.Done()
			localCount := make(map[string]int)

			for _, photo := range part {
				for _, user := range photo.GetUserSimilarPhoto() {
					localCount[user.UserId]++
				}
			}

			// Merge localCount ke global countMap
			mu.Lock()
			for id, cnt := range localCount {
				countMap[id] += cnt
			}
			mu.Unlock()
		}(datas[start:end])
	}

	wg.Wait()
	return countMap
}

// --- Fungsi utama: fetch dari Redis dulu, fallback Postgre ---
func (u *notificationUseCase) fetchFCMTokens(ctx context.Context, userIDs []string) (*[]*entity.UserDevice, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	// 1. Coba ambil dari Redis
	fcmTokens, err := u.redisClient.HMGet(ctx, "fcm_tokens", userIDs...).Result()
	if err != nil {
		return nil, fmt.Errorf("redis HMGET error: %w", err)
	}

	var (
		missingUserIDs []string
		finalResult    []*entity.UserDevice
	)

	for idx, raw := range fcmTokens {
		if raw == nil {
			// Token tidak ditemukan di Redis
			missingUserIDs = append(missingUserIDs, userIDs[idx])
		} else {
			// Token ditemukan di Redis
			tokenStr, ok := raw.(string)
			if !ok {
				log.Printf("Invalid token format for userID: %s", userIDs[idx])
				missingUserIDs = append(missingUserIDs, userIDs[idx])
				continue
			}
			finalResult = append(finalResult, &entity.UserDevice{
				UserId: userIDs[idx],
				Token:  tokenStr,
			})
		}
	}

	// 2. Kalau ada missing, fallback ke PostgreSQL
	if len(missingUserIDs) > 0 {
		dbResults, err := u.userDeviceRepository.FetchFCMTokensFromPostgre(ctx, u.db, missingUserIDs)
		if err != nil {
			return nil, fmt.Errorf("postgres fetch error: %w", err)
		}

		// Update Redis cache untuk yang baru didapat
		pipe := u.redisClient.Pipeline()
		for _, ua := range *dbResults {
			pipe.HSet(ctx, "fcm_tokens", ua.UserId, ua.Token)
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			log.Printf("Redis cache update error: %v", err)
		}

		// Gabungkan hasil dari PostgreSQL ke final result
		finalResult = append(finalResult, *dbResults...)
	}

	log.Print("[FETCHFCM][LEN]", len(finalResult))

	return &finalResult, nil
}

func (u *notificationUseCase) ProcessAndSendNotifications(ctx context.Context, datas []*photopb.BulkUserSimilarPhoto) error {
	countMap := u.countUsersParallel(datas)

	lenCountMap := len(countMap)
	userIDs := make([]string, 0, lenCountMap)
	for idx := range countMap {
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
		u.sendFCMWorkerPool(ctx, userAuthentications, countMap, 10)
	}

	if outerError != nil {
		return outerError
	}
	return nil
}

// --- sendFCMWithRetry function ---
func (u *notificationUseCase) sendFCMWithRetry(ctx context.Context, userID, fcmToken, message string) error {
	const maxRetries = 3
	var backoff = time.Millisecond * 500 // 500ms backoff awal

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("[TOKEN][sendFCMWithRetry]fcm happen trom fcm error %s", fcmToken)
		err := u.sendFCM(ctx, fcmToken, message)

		if err == nil {
			return nil // sukses
		}
		log.Printf("[ERROR][sendFCMWithRetry]Error happen trom fcm error %v", err)
		u.removeUserToken(ctx, userID)

		if helper.IsFCMInvalidTokenError(err) {
			log.Printf("[ERROR][sendFCMWithRetry]Error happen trom fcm error %v", err)
			// Token rusak, hapus
			u.removeUserToken(ctx, userID)
			return fmt.Errorf("invalid token for userID %s: %w", userID, err)
		}

		if helper.IsFCMRetryableError(err) {
			log.Printf("Retryable error sending FCM to userID=%s, attempt=%d, err=%v", userID, attempt, err)
			time.Sleep(backoff)
			backoff *= 2 // exponential backoff
			continue
		}

		// Unknown error, log saja
		log.Printf("Unexpected error sending FCM to userID=%s: %v", userID, err)
		return err
	}

	return fmt.Errorf("failed to send FCM to userID=%s after %d attempts", userID, maxRetries)
}

func (u *notificationUseCase) removeUserToken(ctx context.Context, userID string) (err error) {
	tx, err := repository.BeginTxx(u.db, ctx, u.logs)
	if err != nil {
		return err
	}

	defer func() {
		repository.Rollback(err, tx, ctx, u.logs)
	}()

	// Hapus dari DB
	if err := u.userDeviceRepository.DeleteByUserID(ctx, tx, userID); err != nil {
		return err
	}
	log.Printf("Removing FCM token for user: %s", userID)

	// Hapus dari Redis hash "fcm_tokens"
	if err := u.redisClient.HDel(ctx, "fcm_tokens", userID).Err(); err != nil {
		return fmt.Errorf("redis HDEL error: %w", err)
	}

	// Commit DB transaction
	if err := repository.Commit(tx, u.logs); err != nil {
		return err
	}

	return nil
}

func (u *notificationUseCase) sendFCM(ctx context.Context, fcmToken, message string) error {
	log.Print("[REQUESTED][sendFCM")
	msg := &messaging.Message{
		Token: fcmToken,
		Notification: &messaging.Notification{
			Title: "Foto Mirip Terdeteksi",
			Body:  message, // Contoh: "Terdapat 5 foto yang mirip dengan Anda!"
		},
		Data: map[string]string{
			"type":    "similar_photo",
			"message": message,
		},
	}

	response, err := u.cloudMessagingAdapter.Send(ctx, msg)
	if err != nil {

		if status.Code(err) == codes.NotFound {
			fmt.Println("User not found")
		} else {
			fmt.Println("Error getting user:", err)
		}
		// Type assertion untuk memeriksa apakah error merupakan *googleapi.Error
		if gErr, ok := err.(*googleapi.Error); ok {
			// Menangani error dari Google API secara spesifik
			log.Printf("Google API error - Code: %v, Message: %v", gErr.Code, gErr.Message)

			// Misalnya, jika error 404 (Not Found)
			if gErr.Code == 404 {
				log.Println("Resource not found")
			}
		} else {
			// Jika bukan error dari Google API, cetak error biasa
			log.Printf("Error sending FCM to %s: %v", fcmToken, err)
		}
		return err
	}

	log.Printf("✅ FCM sent to %s: %s", fcmToken, response)
	return nil
}

func (u *notificationUseCase) sendFCMWorkerPool(ctx context.Context, userAuthentications *[]*entity.UserDevice, countMap map[string]int, workerCount int) {
	jobChan := make(chan *entity.UserDevice)
	var wg sync.WaitGroup

	// Start workerCount workers
	for i := range workerCount {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case userAuth, ok := <-jobChan:
					if !ok {
						return
					}

					count := countMap[userAuth.UserId]
					if count == 0 {
						continue
					}

					message := fmt.Sprintf("Terdapat %d foto yang mirip dengan Anda!", count)
					if err := u.sendFCMWithRetry(ctx, userAuth.UserId, userAuth.Token, message); err != nil {
						log.Printf("[Worker %d] Failed send to %s: %v", workerID, userAuth.UserId, err)
					}
				}
			}
		}(i)
	}

loop:
	for _, userAuth := range *userAuthentications {
		select {
		case <-ctx.Done():
			break loop
		case jobChan <- userAuth:
		}
	}

	close(jobChan)

	// Wait all workers finish
	wg.Wait()
}
