package config

import (
	"context"
	"log"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/utils"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	MinioClient     *minio.Client
	minioBucketName string
	enpoint         string
	Logs            *logger.Log
}

func NewMinio() *Minio {
	logger := logger.New("minio")
	ctx := context.Background()
	var minioClient *minio.Client
	var err error

	minioHost := utils.GetEnv("MINIO_HOST")
	minioPort := utils.GetEnv("MINIO_PORT")
	minioRootUser := utils.GetEnv("MINIO_ROOT_USER")
	minioRootPassword := utils.GetEnv("MINIO_ROOT_PASSWORD")
	minioBucket := utils.GetEnv("MINIO_BUCKET")
	minioLocation := utils.GetEnv("MINIO_LOCATION")
	endpoint := minioHost + ":" + minioPort

	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioRootUser, minioRootPassword, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = minioClient.MakeBucket(ctx, minioBucket, minio.MakeBucketOptions{Region: minioLocation})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, minioBucket)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", minioBucket)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", minioBucket)
	}

	log.Printf("Successfully connected %s\n", minioBucket)

	return &Minio{
		MinioClient:     minioClient,
		minioBucketName: minioBucket,
		enpoint:         endpoint,
		Logs:            logger,
	}
}

func (m *Minio) GetBucketName() string {
	BucketName := m.minioBucketName
	return BucketName
}

func (m *Minio) GetEndpoint() string {
	Endpoint := m.enpoint
	return Endpoint
}
