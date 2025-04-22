package adapter

import (
	"be-yourmoments/photo-svc/internal/helper/discovery"
	"context"

	"github.com/be-yourmoments/pb"
)

type AiAdapter interface {
}

type aiAdapter struct {
	client pb.AiServiceClient
}

func NewAiAdapter(ctx context.Context, registry discovery.Registry) (AiAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "ai-svc-grpc", registry)
	if err != nil {
		return nil, err
	}

	client := pb.NewAiServiceClient(conn)

	return &aiAdapter{
		client: client,
	}, nil
}
