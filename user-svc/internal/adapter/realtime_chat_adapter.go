package adapter

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/nullable"
	"github.com/oklog/ulid/v2"
)

type RealtimeChatAdapter interface {
	CreateChatRoom(ctx context.Context, user *entity.User, userProfile *entity.UserProfile)
	GetRoom(ctx context.Context, roomUserId string) ([]*firestore.DocumentSnapshot, error)
	CreateRoom(ctx context.Context, roomUserId string, participants []string) error
	SendMessage(ctx context.Context, roomId, senderId, safeMessage string) error
}

type realtimeChatAdapter struct {
	firestoreClient *firestore.Client
	logs            logger.Log
}

func NewRealtimeChatAdapter(ctx context.Context, app *firebase.App, logs logger.Log) RealtimeChatAdapter {
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing firebase firestore: %v", err)
	}
	return &realtimeChatAdapter{
		firestoreClient: firestoreClient,
		logs:            logs,
	}
}
func (u *realtimeChatAdapter) CreateChatRoom(ctx context.Context, user *entity.User, userProfile *entity.UserProfile) {
	userRef := u.firestoreClient.Collection("users").Doc(user.Id)
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

func (u *realtimeChatAdapter) GetRoom(ctx context.Context, roomUserId string) ([]*firestore.DocumentSnapshot, error) {
	query := u.firestoreClient.
		Collection("rooms").
		Where("roomUserId", "==", roomUserId).
		Limit(1)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func (u *realtimeChatAdapter) CreateRoom(ctx context.Context, roomUserId string, participants []string) error {
	roomId := ulid.Make().String()

	_, err := u.firestoreClient.
		Collection("rooms").
		Doc(roomId).
		Set(ctx, map[string]interface{}{
			"roomId":       roomId,
			"roomUserId":   roomUserId,
			"participants": participants,
			"createdAt":    firestore.ServerTimestamp,
		})

	if err != nil {
		return fmt.Errorf("failed to create a firestore room: %w", err)
	}
	return nil
}

func (u *realtimeChatAdapter) SendMessage(ctx context.Context, roomId, senderId, safeMessage string) error {
	_, _, err := u.firestoreClient.
		Collection("rooms").
		Doc(roomId).
		Collection("messages").
		Add(ctx, map[string]interface{}{
			"senderId":  senderId,
			"message":   safeMessage,
			"timestamp": firestore.ServerTimestamp,
		})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}
