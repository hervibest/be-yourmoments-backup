package v1

// auth := &entity.Auth{
// 	Id:          user.Id,
// 	Username:    user.Username,
// 	Email:       user.Email.String,
// 	PhoneNumber: user.PhoneNumber.String,
// 	CreatorId:   creator.Id,
// 	WalletId:    wallet.Id,
// 	Similarity:  userProfile.Similarity,
// }

// now := time.Now()

// userDevice := &entity.UserDevice{
// 	Id:        ulid.Make().String(),
// 	UserId:    user.Id,
// 	Token:     request.DeviceToken,
// 	Platform:  request.Platform,
// 	CreatedAt: &now,
// }

// _, err = u.userDeviceRepository.Create(ctx, u.db, userDevice)
// if err != nil {
// 	return nil, nil, helper.WrapInternalServerError(u.logs, "failed to create user device", err)
// }

// setKey := fmt.Sprintf("fcm_tokens:%s", user.Id)
// if err := u.cacheAdapter.SAdd(ctx, setKey, request.DeviceToken); err != nil {
// 	return nil, nil, helper.WrapInternalServerError(u.logs, "failed to SAdd redis set", err)
// }

// token, err := u.generateToken(ctx, auth)
// if err != nil {
// 	return nil, nil, helper.WrapInternalServerError(u.logs, "failed find user by email not google", err)
// }
