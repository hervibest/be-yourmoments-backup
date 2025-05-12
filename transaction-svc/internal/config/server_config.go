package config

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/utils"
)

var EndpointPrefix = utils.GetEnv("ENDPOINT_PREFIX")

type ServerConfig struct {
	HTTP             string
	HTTPAddr         string
	HTTPInternalAddr string
	HTTPPort         string
	GRPC             string
	GRPCAddr         string
	GRPCInternalAddr string
	GRPCPort         string
	ConsulAddr       string
	Name             string
}

func NewServerConfig() ServerConfig {
	httpAddr := utils.GetEnv("HTTP_ADDR")
	if httpAddr == "" {
		log.Fatal("HTTP_ADDR environment variable is not set")
	}
	// Ambil IP container dari hostname
	hostname, _ := os.Hostname()
	addrs, _ := net.LookupIP(hostname)

	var ip string
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil && !ipv4.IsLoopback() {
			ip = ipv4.String()
			break
		}
	}
	httpInternalAddr := ip
	port := utils.GetEnv("HTTP_PORT")
	if port == "" {
		log.Fatal("HTTP_PORT environment variable is not set")
	}
	grpcAddr := utils.GetEnv("GRPC_ADDR")
	if grpcAddr == "" {
		log.Fatal("GRPC_ADDR environment variable is not set")
	}
	grpcInternalAddr := ip
	grpcPort := utils.GetEnv("GRPC_PORT")
	if grpcPort == "" {
		log.Fatal("GRPC_PORT environment variable is not set")
	}
	consulAddr := utils.GetEnv("CONSUL_ADDR")
	if consulAddr == "" {
		log.Fatal("CONSUL_ADDR environment variable is not set")
	}
	name := utils.GetEnv("SERVICE_NAME")
	if name == "" {
		log.Fatal("SERVICE_NAME environment variable is not set")
	}
	return ServerConfig{
		HTTP:             fmt.Sprintf("%s:%s", httpAddr, port),
		HTTPAddr:         httpAddr,
		HTTPInternalAddr: httpInternalAddr,
		HTTPPort:         port,
		GRPC:             fmt.Sprintf("%s:%s", grpcAddr, grpcPort),
		GRPCAddr:         grpcAddr,
		GRPCInternalAddr: grpcInternalAddr,
		GRPCPort:         grpcPort,
		ConsulAddr:       consulAddr,
		Name:             name,
	}
}
