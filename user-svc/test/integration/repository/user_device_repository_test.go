package integrationrepo

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/config"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/repository"
	"github.com/oklog/ulid/v2"
)

func TestCreateDeviceRepo(t *testing.T) {
	dbConfig := config.NewDB()
	deviceRepository := repository.NewUserDeviceRepository()

	now := time.Now()
	userDevice := &entity.UserDevice{
		Id:        ulid.Make().String(),
		UserId:    "01JSH02TEM8ZJ5GRWJH04XECH2",
		Token:     "",
		Platform:  "WEB",
		CreatedAt: &now,
	}

	_, err := deviceRepository.Create(context.Background(), dbConfig, userDevice)
	if err != nil {
		log.Default().Printf("Error happen when create device repo %v", err)
	}
}
