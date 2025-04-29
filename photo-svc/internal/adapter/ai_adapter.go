package adapter

import (
	"be-yourmoments/photo-svc/internal/helper/discovery"
	"context"

	aipb "github.com/be-yourmoments/pb/ai"
)

type AiAdapter interface {
}

type aiAdapter struct {
	client aipb.AiServiceClient
}

func NewAiAdapter(ctx context.Context, registry discovery.Registry) (AiAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "ai-svc-grpc", registry)
	if err != nil {
		return nil, err
	}

	client := aipb.NewAiServiceClient(conn)

	return &aiAdapter{
		client: client,
	}, nil
}
