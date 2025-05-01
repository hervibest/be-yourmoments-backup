package v1

/* USING REDIS HMAP (HASH MAP) */
// // --- Fungsi utama: fetch dari Redis dulu, fallback Postgre ---
// func (u *notificationUseCase) fetchFCMTokens(ctx context.Context, userIDs []string) (*[]*entity.UserDevice, error) {
// 	if len(userIDs) == 0 {
// 		return nil, nil
// 	}

// 	// 1. Coba ambil dari Redis
// 	fcmTokens, err := u.redisClient.HMGet(ctx, "fcm_tokens", userIDs...).Result()
// 	if err != nil {
// 		return nil, fmt.Errorf("redis HMGET error: %w", err)
// 	}

// 	var (
// 		missingUserIDs []string
// 		finalResult    []*entity.UserDevice
// 	)

// 	for idx, raw := range fcmTokens {
// 		if raw == nil {
// 			// Token tidak ditemukan di Redis
// 			missingUserIDs = append(missingUserIDs, userIDs[idx])
// 		} else {
// 			// Token ditemukan di Redis
// 			tokenStr, ok := raw.(string)
// 			if !ok {
// 				log.Printf("Invalid token format for userID: %s", userIDs[idx])
// 				missingUserIDs = append(missingUserIDs, userIDs[idx])
// 				continue
// 			}
// 			finalResult = append(finalResult, &entity.UserDevice{
// 				UserId: userIDs[idx],
// 				Token:  tokenStr,
// 			})
// 		}
// 	}

// 	// 2. Kalau ada missing, fallback ke PostgreSQL
// 	if len(missingUserIDs) > 0 {
// 		dbResults, err := u.userDeviceRepository.FetchFCMTokensFromPostgre(ctx, u.db, missingUserIDs)
// 		if err != nil {
// 			return nil, fmt.Errorf("postgres fetch error: %w", err)
// 		}

// 		// Update Redis cache untuk yang baru didapat
// 		pipe := u.redisClient.Pipeline()
// 		for _, ua := range *dbResults {
// 			pipe.HSet(ctx, "fcm_tokens", ua.UserId, ua.Token)
// 		}
// 		_, err = pipe.Exec(ctx)
// 		if err != nil {
// 			log.Printf("Redis cache update error: %v", err)
// 		}

// 		// Gabungkan hasil dari PostgreSQL ke final result
// 		finalResult = append(finalResult, *dbResults...)
// 	}

// 	log.Print("[FETCHFCM][LEN]", len(finalResult))

// 	return &finalResult, nil
// }

// SEND SINGLE FCM
// func (u *notificationUseCase) sendFCMWithRetry(ctx context.Context, userID, fcmToken, message string) error {
// 	const maxRetries = 3
// 	backoff := 500 * time.Millisecond

// 	for attempt := 1; attempt <= maxRetries; attempt++ {
// 		log.Printf("[SEND FCM] Attempt %d: Sending to token=%s", attempt, fcmToken)

// 		err := u.sendFCM(ctx, fcmToken, message)
// 		if err == nil {
// 			log.Printf("[SEND FCM] Success on attempt %d for userID=%s", attempt, userID)
// 			return nil
// 		}

// 		fcmErr := helper.ParseFCMError(err)

// 		switch {
// 		case fcmErr.IsInvalidToken():
// 			log.Printf("[INVALID TOKEN] userID=%s, token=%s, code=%s", userID, fcmToken, fcmErr.Code)

// 			if removeErr := u.removeUserToken(ctx, userID, fcmToken); removeErr != nil {
// 				log.Printf("[REMOVE TOKEN ERROR] userID=%s, err=%v", userID, removeErr)
// 			}

// 			return fmt.Errorf("invalid FCM token for userID=%s: %w", userID, err)

// 		case fcmErr.IsRetryable():
// 			log.Printf("[RETRYABLE] userID=%s, token=%s, code=%s → retrying after %v", userID, fcmToken, fcmErr.Code, backoff)
// 			time.Sleep(backoff)
// 			backoff *= 2
// 			continue

// 		case fcmErr.IsAuthError():
// 			log.Printf("[AUTH ERROR] userID=%s, token=%s, code=%s → admin action required", userID, fcmToken, fcmErr.Code)
// 			u.alertAdminFCMAuthIssue(userID, fcmErr) // Anda buat fungsi ini sesuai kebutuhan
// 			return fmt.Errorf("auth error sending FCM to userID=%s: %w", userID, err)

// 		default:
// 			log.Printf("[UNEXPECTED FCM ERROR] userID=%s, token=%s, code=%s, details=%s", userID, fcmToken, fcmErr.Code, fcmErr.Details)
// 			return fmt.Errorf("unexpected FCM error: %w", err)
// 		}
// 	}

// 	return fmt.Errorf("exceeded max retries sending FCM to userID=%s, token=%s", userID, fcmToken)
// }

// TODO TYPE ASSERTION
// USECASE BAGAIMAN KALAU FCM DOWN BESAR BESAR RAN (RETRY)
// func (u *notificationUseCase) sendFCM(ctx context.Context, fcmToken, message string) error {
// 	log.Print("[REQUESTED][sendFCM")
// 	msg := &messaging.Message{
// 		Token: fcmToken,
// 		Notification: &messaging.Notification{
// 			Title: "Foto Mirip Terdeteksi",
// 			Body:  message, // Contoh: "Terdapat 5 foto yang mirip dengan Anda!"
// 		},
// 		Data: map[string]string{
// 			"type":    "similar_photo",
// 			"message": message,
// 		},
// 	}

// 	response, err := u.cloudMessagingAdapter.Send(ctx, msg)
// 	if err != nil {
// 		return err
// 	}

// 	log.Printf("✅ FCM sent to %s: %s", fcmToken, response)
// 	return nil
// }

// func (u *notificationUseCase) sendFCMWorkerPool(ctx context.Context, userAuthentications *[]*entity.UserDevice, countMap map[string]int, workerCount int) {
// 	jobChan := make(chan *entity.UserDevice)
// 	var wg sync.WaitGroup

// 	// Start workerCount workers
// 	for i := range workerCount {
// 		wg.Add(1)
// 		go func(workerID int) {
// 			defer wg.Done()
// 			for {
// 				select {
// 				case <-ctx.Done():
// 					return
// 				case userAuth, ok := <-jobChan:
// 					if !ok {
// 						return
// 					}

// 					count := countMap[userAuth.UserId]
// 					if count == 0 {
// 						continue
// 					}

// 					message := fmt.Sprintf("Terdapat %d foto yang mirip dengan Anda!", count)
// 					if err := u.sendFCMWithRetry(ctx, userAuth.UserId, userAuth.Token, message); err != nil {
// 						log.Printf("[Worker %d] Failed send to %s: %v", workerID, userAuth.UserId, err)
// 					}
// 				}
// 			}
// 		}(i)
// 	}

// loop:
// 	for _, userAuth := range *userAuthentications {
// 		select {
// 		case <-ctx.Done():
// 			break loop
// 		case jobChan <- userAuth:
// 		}
// 	}

// 	close(jobChan)

// 	// Wait all workers finish
// 	wg.Wait()
// }
