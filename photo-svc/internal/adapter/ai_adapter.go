package adapter

import (
	"context"
	"log"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/discovery"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/utils"

	aipb "github.com/hervibest/be-yourmoments-backup/pb/ai"
)

type AiAdapter interface {
}

type aiAdapter struct {
	client aipb.AiServiceClient
}

func NewAiAdapter(ctx context.Context, registry discovery.Registry, logs *logger.Log) (AiAdapter, error) {
	aiServiceName := utils.GetEnv("AI_SVC_NAME")
	conn, err := discovery.NewGrpcClient(aiServiceName)
	if err != nil {
		logs.CustomError("failed to connect to ai service due to an error : ", err)
		return nil, err
	}

	log.Printf("successfuly connected to %s", aiServiceName)
	client := aipb.NewAiServiceClient(conn)

	return &aiAdapter{
		client: client,
	}, nil
}
