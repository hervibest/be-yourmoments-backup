package uploadconsumer

import (
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/usecase"
	"github.com/nats-io/nats.go"
)

// consumer.go
type UploadConsumer struct {
	photoWorkerUC    usecase.PhotoWorkerUseCase
	facecameWorkerUC usecase.FacecamUseCaseWorker
	js               nats.JetStreamContext
	logs             *logger.Log
	subjects         []string
	durableNames     map[string]string
}

func NewUploadConsumer(
	photoWorkerUC usecase.PhotoWorkerUseCase,
	facecamWorkerUC usecase.FacecamUseCaseWorker,
	js nats.JetStreamContext,
	logs *logger.Log,
) *UploadConsumer {
	return &UploadConsumer{
		photoWorkerUC:    photoWorkerUC,
		facecameWorkerUC: facecamWorkerUC,
		js:               js,
		logs:             logs,
		subjects: []string{
			"upload.bulk.photo",
			"upload.single.facecam",
			"upload.single.photo",
			"upload.update.photo",
		},
		durableNames: map[string]string{
			"upload.bulk.photo":     "upload_bulk_photo_consumer",
			"upload.single.facecam": "upload_single_facecam_consumer",
			"upload.single.photo":   "upload_single_photo_consumer",
			"upload.update.photo":   "upload_update_photo_consumer",
		},
	}

}
