package consul

import (
	"context"
	"fmt"
	"time"

	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

func NewConsulResolverBuilder(client *consul.Client) resolver.Builder {
	return &consulResolverBuilder{
		client: client,
	}
}

type consulResolverBuilder struct {
	client *consul.Client
}

func (b *consulResolverBuilder) Scheme() string {
	return "consul"
}
func (b *consulResolverBuilder) Build(
	target resolver.Target,
	cc resolver.ClientConn,
	opts resolver.BuildOptions,
) (resolver.Resolver, error) {

	ctx, cancel := context.WithCancel(context.Background())

	r := &consulResolver{
		client:      b.client,
		serviceName: target.Endpoint(),
		cc:          cc,
		ctx:         ctx,
		cancel:      cancel,
		logger: func(format string, v ...any) {
			fmt.Printf("[CONSUL-RESOLVER] "+format+"\n", v...)
		},
	}

	r.logger("build resolver for service=%s", r.serviceName)

	go r.watch()
	return r, nil
}

type consulResolver struct {
	client      *consul.Client
	serviceName string
	cc          resolver.ClientConn
	ctx         context.Context
	cancel      context.CancelFunc
	logger      func(format string, v ...any)
}

func (r *consulResolver) ResolveNow(o resolver.ResolveNowOptions) {
	r.update()
}

func (r *consulResolver) Close() {
	if r.cancel != nil {
		r.cancel()
	}
}

func (r *consulResolver) watch() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	r.logger("start watching service=%s", r.serviceName)

	for {
		select {
		case <-ticker.C:
			r.update()

		case <-r.ctx.Done():
			r.logger("resolver stopped")
			return
		}
	}
}

func (r *consulResolver) update() {
	services, _, err := r.client.Health().Service(
		r.serviceName,
		"",
		true, // only passing
		nil,
	)
	if err != nil {
		r.logger("health check error: %v", err)
		return
	}

	if len(services) == 0 {
		r.logger("no healthy instance found")
		// tetap update state kosong → trigger reconnect
		r.cc.UpdateState(resolver.State{
			Addresses: nil,
		})
		return
	}

	addrs := make([]resolver.Address, 0, len(services))
	for _, svc := range services {
		addr := fmt.Sprintf(
			"%s:%d",
			svc.Service.Address,
			svc.Service.Port,
		)
		addrs = append(addrs, resolver.Address{Addr: addr})
	}

	r.logger("update address list: %v", addrs)

	r.cc.UpdateState(resolver.State{
		Addresses: addrs,
	})
}
