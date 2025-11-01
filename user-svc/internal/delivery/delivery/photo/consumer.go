package consumer

import (
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/usecase"
	"github.com/nats-io/nats.go"
)

// consumer.go
type PhotoConsumer struct {
	userUseCase  usecase.UserUseCase
	js           nats.JetStreamContext
	logs         logger.Log
	subjects     []string
	durableNames map[string]string
}

func NewPhotoConsumer(
	userUseCase usecase.UserUseCase,
	js nats.JetStreamContext,
	logs logger.Log,
) *PhotoConsumer {
	return &PhotoConsumer{
		userUseCase: userUseCase,
		js:          js,
		logs:        logs,
		subjects: []string{
			"photo.persist.facecam",
		},
		durableNames: map[string]string{
			"photo.persist.facecam": "photo_persist_facecam_consumer",
		},
	}
}
