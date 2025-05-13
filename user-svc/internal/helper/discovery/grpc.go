package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Registry interface {
	RegisterService(ctx context.Context, serviceName, serviceID, serviceAddress string, servicePort int, tags []string) error
	DeregisterService(ctx context.Context, serviceID string) error
	GetService(ctx context.Context, serviceName string) ([]*consul.ServiceEntry, error)
	HealthCheck(serviceID, serviceName string) error
}

func GenerateServiceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}

func ServiceConnection(ctx context.Context, serviceName string, registry Registry, logs logger.Log) (*grpc.ClientConn, error) {
	const (
		maxRetries = 3
		retryDelay = 10 * time.Second
	)

	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		logs.Log(fmt.Sprintf("trying to connect service: %s with attempt: %d and max retries: %d", serviceName, attempt, maxRetries))
		service, err := registry.GetService(ctx, serviceName)
		if err != nil {
			lastErr = fmt.Errorf("failed to get service: %w", err)
			time.Sleep(retryDelay)
			continue
		}

		if len(service) == 0 {
			lastErr = fmt.Errorf("service %s not found", serviceName)
			time.Sleep(retryDelay)
			continue
		}

		serviceEntry := service[rand.Intn(len(service))]

		conn, err := grpc.NewClient(
			fmt.Sprintf("%s:%d", serviceEntry.Service.Address, serviceEntry.Service.Port),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			lastErr = fmt.Errorf("failed to connect: %w", err)
			time.Sleep(retryDelay)
			continue
		}

		// Berhasil terkoneksi
		return conn, nil
	}

	// Setelah semua percobaan gagal
	return nil, fmt.Errorf("service connection failed after %d retries: %w", maxRetries, lastErr)
}
