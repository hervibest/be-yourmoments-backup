package adapter

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/config"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/model"

	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
)

type StorageAdapter interface {
	UploadCompressedFile(ctx context.Context, file *multipart.FileHeader, uploadFile multipart.File, path string) (*model.MinioFileResponse, error)
	DeleteFile(ctx context.Context, fileName string) (bool, error)
	UploadOriginalFileWithoutMultipart(ctx context.Context, file *multipart.FileHeader, uploadFile io.Reader, path, photoID string) (*model.MinioFileResponse, error)
}

type storageAdapter struct {
	minio *config.Minio
}

func NewStorageAdapter(minio *config.Minio) StorageAdapter {
	return &storageAdapter{
		minio: minio,
	}
}
func (a *storageAdapter) UploadCompressedFile(ctx context.Context, file *multipart.FileHeader, uploadFile multipart.File, path string) (*model.MinioFileResponse, error) {
	// ulidStr := ulid.Make().String()
	// ext := filepath.Ext(file.Filename) // e.g., .jpg
	// cleanFilename := strings.TrimSuffix(file.Filename, ext)
	// safeFilename := strings.ReplaceAll(cleanFilename, " ", "_")

	fileKey := fmt.Sprintf("%s/%s", path, file.Filename)
	contentType := file.Header.Get("Content-Type")

	return a.uploadToMinIO(ctx, uploadFile, fileKey, contentType, file.Filename, file.Size)
}

func (a *storageAdapter) UploadOriginalFileWithoutMultipart(ctx context.Context, file *multipart.FileHeader, uploadFile io.Reader, path, photoID string) (*model.MinioFileResponse, error) {
	ext := filepath.Ext(file.Filename)
	cleanFilename := strings.TrimSuffix(file.Filename, ext)
	safeFilename := strings.ReplaceAll(cleanFilename, " ", "_")

	fileKey := fmt.Sprintf("%s/%s_%s%s", path, safeFilename, photoID, ext)
	contentType := file.Header.Get("Content-Type")

	return a.uploadToMinIO(ctx, uploadFile, fileKey, contentType, file.Filename, file.Size)
}

func (a *storageAdapter) uploadToMinIO(ctx context.Context, uploadFile io.Reader, fileKey, contentType, originalFilename string, size int64) (*model.MinioFileResponse, error) {
	s3PutObjectOutput, err := a.minio.MinioClient.PutObject(ctx, a.minio.GetBucketName(), fileKey, uploadFile, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to put object to minio storage : %v", err)
	}

	fileURL, err := a.minio.MinioClient.PresignedGetObject(ctx, a.minio.GetBucketName(), fileKey, 1*time.Hour, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL : %v", err)
	}

	return &model.MinioFileResponse{
		ChecksumCRC32:  s3PutObjectOutput.ChecksumCRC32,
		ChecksumCRC32C: s3PutObjectOutput.ChecksumCRC32C,
		ChecksumSHA1:   s3PutObjectOutput.ChecksumSHA1,
		ChecksumSHA256: s3PutObjectOutput.ChecksumSHA256,
		ETag:           s3PutObjectOutput.ETag,
		Expiration:     s3PutObjectOutput.Expiration,
		URL:            fileURL.String(),
		Filename:       originalFilename,
		FileKey:        fileKey,
		Mimetype:       contentType,
		Size:           size,
	}, nil
}

func (a *storageAdapter) DeleteFile(ctx context.Context, fileName string) (bool, error) {
	err := a.minio.MinioClient.RemoveObject(ctx, a.minio.GetBucketName(), fileName, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		return false, fmt.Errorf("failed to delete image in minio storage : %v", err)
	}

	return true, nil
}
