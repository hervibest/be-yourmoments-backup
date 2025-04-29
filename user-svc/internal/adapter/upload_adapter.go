package adapter

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/config"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

type UploadAdapter interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, uploadFile multipart.File, path string) (*model.MinioFileResponse, error)
	DeleteFile(ctx context.Context, fileName string) (bool, error)
	GetPresignedUrl(ctx context.Context, fileKey string) (string, error)
}

type uploadAdapter struct {
	minio       *config.Minio
	redisClient *redis.Client // Redis client
}

func NewUploadAdapter(minio *config.Minio, redisClient *redis.Client) UploadAdapter {
	return &uploadAdapter{
		minio:       minio,
		redisClient: redisClient,
	}
}

func (a *uploadAdapter) UploadFile(ctx context.Context, file *multipart.FileHeader, uploadFile multipart.File, path string) (*model.MinioFileResponse, error) {
	fileKey := path + string(RandomNumber(31)) + "_" + file.Filename
	contentType := file.Header.Get("Content-Type")

	s3PutObjectOutput, err := a.minio.MinioClient.PutObject(ctx, a.minio.GetBucketName(), fileKey, uploadFile, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		a.minio.Logs.Error("failed to upload file to S3" + err.Error())

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fileResponse := new(model.MinioFileResponse)
	fileResponse.ChecksumCRC32 = s3PutObjectOutput.ChecksumCRC32
	fileResponse.ChecksumCRC32C = s3PutObjectOutput.ChecksumCRC32C
	fileResponse.ChecksumSHA1 = s3PutObjectOutput.ChecksumSHA1
	fileResponse.ChecksumSHA256 = s3PutObjectOutput.ChecksumSHA256
	fileResponse.ETag = s3PutObjectOutput.ETag
	fileResponse.Expiration = s3PutObjectOutput.Expiration

	fileURL, err := a.minio.MinioClient.PresignedGetObject(ctx, a.minio.GetBucketName(), fileKey, 1*time.Hour, nil)
	if err != nil {
		a.minio.Logs.Error("failed to generate presigned URL:" + err.Error())

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fileResponse.URL = fileURL.String()
	fileResponse.Filename = file.Filename
	fileResponse.FileKey = fileKey
	fileResponse.Mimetype = contentType
	fileResponse.Size = file.Size

	return fileResponse, nil
}

func (a *uploadAdapter) DeleteFile(ctx context.Context, fileName string) (bool, error) {

	err := a.minio.MinioClient.RemoveObject(ctx, a.minio.GetBucketName(), fileName, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		return false, fmt.Errorf("failed to delete file: %w", err)
	}

	return true, nil
}

func (a *uploadAdapter) GetPresignedUrl(ctx context.Context, fileKey string) (string, error) {
	cacheKey := "presigned_url:" + fileKey

	if url, err := a.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		return url, nil
	}

	url, err := a.minio.MinioClient.PresignedGetObject(ctx, a.minio.GetBucketName(), fileKey, 15*time.Minute, nil)
	if err != nil {
		return "", fmt.Errorf("minio client error presigned object : %+v", err)
	}

	urlStr := url.String()

	err = a.redisClient.Set(ctx, cacheKey, urlStr, 15*time.Minute).Err()
	if err != nil {
		return "", fmt.Errorf("redis set cached presigned url : %+v", err)
	}

	return urlStr, nil
}
