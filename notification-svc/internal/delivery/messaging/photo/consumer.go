package consumer

import (
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/usecase"
	"github.com/nats-io/nats.go"
)

// consumer.go
type PhotoConsumer struct {
	notficiationUseCase usecase.NotificationUseCase
	js                  nats.JetStreamContext
	logs                logger.Log
	subjects            []string
	durableNames        map[string]string
}

func NewPhotoConsumer(
	notficiationUseCase usecase.NotificationUseCase,
	js nats.JetStreamContext,
	logs logger.Log,
) *PhotoConsumer {
	return &PhotoConsumer{
		notficiationUseCase: notficiationUseCase,
		js:                  js,
		logs:                logs,
		subjects: []string{
			"photo.bulk",
			"photo.single.facecam",
			"photo.single.photo",
		},
		durableNames: map[string]string{
			"photo.bulk":           "photo_bulk_consumer",
			"photo.single.facecam": "photo_single_facecam_consumer",
			"photo.single.photo":   "photo_single_photo_consumer",
		},
	}
}
