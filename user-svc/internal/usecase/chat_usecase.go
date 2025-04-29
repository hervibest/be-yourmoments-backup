package usecase

import (
	"be-yourmoments/user-svc/internal/adapter"
	errorcode "be-yourmoments/user-svc/internal/enum/error"
	"be-yourmoments/user-svc/internal/helper"
	"be-yourmoments/user-svc/internal/helper/logger"
	"be-yourmoments/user-svc/internal/model"
	"context"
	"html"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/oklog/ulid/v2"
)

type ChatUseCase interface {
	GetCustomToken(ctx context.Context, req *model.RequestCustomToken) (*model.CustomTokenResponse, error)
	GetOrCreateRoom(ctx context.Context, req *model.RequestGetOrCreateRoom) (*model.GetOrCreateRoomResponse, error)
	SendMessage(ctx context.Context, req *model.RequestSendMessage) error
}

type chatUseCase struct {
	firestoreClientAdapter adapter.FirestoreClientAdapter
	authClientAdapter      adapter.AuthClientAdapter
	cloudMessagingAdapter  adapter.CloudMessagingAdapter
	perspectiveAdapter     adapter.PerspectiveAdapter
	logs                   *logger.Log
}

func NewChatUseCase(firestoreClientAdapter adapter.FirestoreClientAdapter, authClientAdapter adapter.AuthClientAdapter,
	cloudMessagingAdapter adapter.CloudMessagingAdapter, perspectiveAdapter adapter.PerspectiveAdapter, logs *logger.Log) ChatUseCase {
	return &chatUseCase{
		firestoreClientAdapter: firestoreClientAdapter,
		authClientAdapter:      authClientAdapter,
		cloudMessagingAdapter:  cloudMessagingAdapter,
		perspectiveAdapter:     perspectiveAdapter,
		logs:                   logs,
	}
}

func (u *chatUseCase) GetOrCreateRoom(ctx context.Context, req *model.RequestGetOrCreateRoom) (*model.GetOrCreateRoomResponse, error) {
	roomUserId := generateRoomId(req.SenderId, req.ReceiverId)
	created := false
	var roomId string

	query := u.firestoreClientAdapter.
		Collection("rooms").
		Where("roomUserId", "==", roomUserId).
		Limit(1)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to query firestore : ", err)
	}

	if len(docs) > 0 {
		data := docs[0].Data()

		if rn, ok := data["roomId"].(string); ok {
			roomId = rn
		} else {
			return nil, helper.WrapInternalServerError(u.logs, "failed to get roomId : roomId missing in documment : ", err)
		}
	} else {
		roomId = ulid.Make().String()
		participants := []string{req.SenderId, req.ReceiverId}

		_, err := u.firestoreClientAdapter.
			Collection("rooms").
			Doc(roomId).
			Set(ctx, map[string]interface{}{
				"roomId":       roomId,
				"roomUserId":   roomUserId,
				"participants": participants,
				"createdAt":    firestore.ServerTimestamp,
			})
		if err != nil {
			return nil, helper.WrapInternalServerError(u.logs, "failed to create a firestore room : ", err)
		}

		created = true
	}

	return &model.GetOrCreateRoomResponse{
		RoomId:  roomId,
		Created: created,
	}, nil
}

func generateRoomId(userA, userB string) string {
	if userA < userB {
		return userA + "_" + userB
	}
	return userB + "_" + userA
}

func (u *chatUseCase) GetCustomToken(ctx context.Context, req *model.RequestCustomToken) (*model.CustomTokenResponse, error) {
	token, err := u.authClientAdapter.CustomToken(ctx, req.UserId)
	if err != nil {
		return nil, helper.WrapInternalServerError(u.logs, "failed to check toxic from perpsective api : ", err)
	}

	response := &model.CustomTokenResponse{
		Token: token,
	}

	return response, nil
}

func (u *chatUseCase) SendMessage(ctx context.Context, req *model.RequestSendMessage) error {
	trimmed := strings.TrimSpace(req.Message)
	if trimmed == "" {
		return helper.NewUseCaseError(errorcode.ErrValidationFailed, "Empty message not allowed")
	}

	if len(trimmed) > 500 {
		return helper.NewUseCaseError(errorcode.ErrValidationFailed, "Message too long. Make sure only 500 words")
	}

	isToxic, err := u.perspectiveAdapter.IsToxicMessage(trimmed)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to check toxic from perpsective api", err)
	}

	if isToxic {
		return helper.NewUseCaseError(errorcode.ErrValidationFailed, "Toxic content detected")
	}

	safeMessage := html.EscapeString(trimmed)
	_, _, err = u.firestoreClientAdapter.
		Collection("rooms").
		Doc(req.RoomId).
		Collection("messages").
		Add(ctx, map[string]interface{}{
			"senderId":  req.SenderId,
			"message":   safeMessage,
			"timestamp": firestore.ServerTimestamp,
		})

	if err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to send message to firestore : ", err)
	}

	return nil
}
