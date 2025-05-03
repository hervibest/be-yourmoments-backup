package adapter

import (
	"io"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/config"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"

	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

type StorageAdapter interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, uploadFile multipart.File, path string) (*model.MinioFileResponse, error)
	DeleteFile(ctx context.Context, fileName string) (bool, error)
	GetFile(ctx context.Context, fileName string) (io.ReadCloser, error)
}

type storageAdapter struct {
	minio *config.Minio
}

func NewStorageAdapter(minio *config.Minio) StorageAdapter {
	return &storageAdapter{
		minio: minio,
	}
}

func (a *storageAdapter) UploadFile(ctx context.Context, file *multipart.FileHeader, uploadFile multipart.File, path string) (*model.MinioFileResponse, error) {
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

func (a *storageAdapter) DeleteFile(ctx context.Context, fileName string) (bool, error) {
	err := a.minio.MinioClient.RemoveObject(ctx, a.minio.GetBucketName(), fileName, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		return false, fmt.Errorf("failed to delete file: %w", err)
	}

	return true, nil
}

func (a *storageAdapter) GetFile(ctx context.Context, fileName string) (io.ReadCloser, error) {
	object, err := a.minio.MinioClient.GetObject(ctx, a.minio.GetBucketName(), fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	_, err = object.Stat()
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, fmt.Errorf("file not found: %w", err)
		}
		return nil, fmt.Errorf("failed to stat object: %w", err)
	}

	return object, nil
}
