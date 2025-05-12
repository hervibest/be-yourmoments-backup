package adapter

import (
	"context"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/discovery"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/utils"

	aipb "github.com/hervibest/be-yourmoments-backup/pb/ai"
)

type AiAdapter interface {
}

type aiAdapter struct {
	client aipb.AiServiceClient
}

func NewAiAdapter(ctx context.Context, registry discovery.Registry) (AiAdapter, error) {
	aiServiceName := utils.GetEnv("AI_SVC_NAME")
	conn, err := discovery.ServiceConnection(ctx, aiServiceName, registry)
	if err != nil {
		return nil, err
	}

	client := aipb.NewAiServiceClient(conn)

	return &aiAdapter{
		client: client,
	}, nil
}
