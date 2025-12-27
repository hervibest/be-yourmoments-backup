package discovery

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGrpcClient(serviceName string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("consul:///%s", serviceName), // ⬅️ TIGA SLASH
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(retryDebugInterceptor()),
		grpc.WithDefaultServiceConfig(`{
		"loadBalancingPolicy": "round_robin",
		"retryPolicy": {
			"maxAttempts": 3,
			"initialBackoff": "0.1s",
			"maxBackoff": "1s",
			"backoffMultiplier": 2,
			"retryableStatusCodes": ["UNAVAILABLE"]
		}
	}`),
	)
	if err != nil {
		return nil, err
	}
	return conn, err
}

func retryDebugInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		if err != nil {
			log.Printf(
				"[gRPC-CLIENT] method=%s error=%v elapsed=%s",
				method,
				err,
				time.Since(start),
			)
		}

		return err
	}
}
