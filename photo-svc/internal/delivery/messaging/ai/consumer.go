package aiconsumer

import (
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/usecase"
	"github.com/nats-io/nats.go"
)

// consumer.go
type AIConsumer struct {
	userSimilarWorkerUC usecase.UserSimilarWorkerUseCase
	js                  nats.JetStreamContext
	logs                *logger.Log
	subjects            []string
	durableNames        map[string]string
}

func NewAIConsumer(
	userSimilarWorkerUC usecase.UserSimilarWorkerUseCase,
	js nats.JetStreamContext,
	logs *logger.Log,
) *AIConsumer {
	return &AIConsumer{
		userSimilarWorkerUC: userSimilarWorkerUC,
		js:                  js,
		logs:                logs,
		subjects: []string{
			"ai.bulk.photo",
			"ai.single.facecam",
			"ai.single.photo",
		},
		durableNames: map[string]string{
			"ai.bulk.photo":     "ai_bulk_consumer",
			"ai.single.facecam": "ai_single_facecam_consumer",
			"ai.single.photo":   "ai_single_ai_consumer",
		},
	}
}
