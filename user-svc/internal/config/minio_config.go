package config

import (
	"context"
	"log"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/utils"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	MinioClient           *minio.Client
	minioBucketName       string
	minioPublicBucketName string
	enpoint               string
	Logs                  logger.Log
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
	minioPublicBucket := utils.GetEnv("MINIO_PUBLIC_BUCKET")
	minioLocation := utils.GetEnv("MINIO_LOCATION")
	endpoint := minioHost + ":" + minioPort

	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioRootUser, minioRootPassword, ""),
		Secure: false,
	})
	logger.CustomError("failed to connect to minio", err)

	buckets := []string{minioBucket, minioPublicBucket}

	for _, bucket := range buckets {
		err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: minioLocation})
		if err != nil {
			exists, errBucketExists := minioClient.BucketExists(ctx, bucket)
			if errBucketExists == nil && exists {
				log.Printf("We already own %s\n", bucket)
			} else {
				log.Fatalln("Minio config error", err)
			}
		} else {
			log.Printf("Successfully created %s\n", bucket)
		}
		log.Printf("Successfully connected %s\n", bucket)
	}

	return &Minio{
		MinioClient:           minioClient,
		minioBucketName:       minioBucket,
		minioPublicBucketName: minioPublicBucket,
		enpoint:               endpoint,
		Logs:                  logger,
	}
}

func (m *Minio) GetBucketName() string {
	return m.minioBucketName
}

func (m *Minio) GetPublicBucketName() string {
	return m.minioPublicBucketName
}

func (m *Minio) GetPublicURL() string {
	return m.enpoint + "/" + m.minioPublicBucketName
}

func (m *Minio) GetEndpoint() string {
	return m.enpoint
}
