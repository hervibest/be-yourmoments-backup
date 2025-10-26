package config

import (
	"log"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/utils"
)

const (
	EnvLocal       = "local"
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

var environtment string

func init() {
	environtment = utils.GetEnv("ENVIRONMENT")
}

func Get() string {
	log.Default().Println("Current Environment: " + environtment)
	return environtment
}

func IsLocal() bool {
	return Get() == EnvLocal
}

func IsDevelopment() bool {
	return Get() == EnvDevelopment
}

func IsProduction() bool {
	return Get() == EnvProduction
}
