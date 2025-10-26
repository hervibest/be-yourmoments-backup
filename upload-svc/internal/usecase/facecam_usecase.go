package usecase

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/entity"
	errorcode "github.com/hervibest/be-yourmoments-backup/upload-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/logger"
	producer "github.com/hervibest/be-yourmoments-backup/upload-svc/internal/messaging"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/model"

	"github.com/oklog/ulid/v2"
)

type FacecamUseCase interface {
	UploadFacecam(ctx context.Context, file *multipart.FileHeader, userId, creatorId string) error
	// UpdateProcessedPhoto(ctx context.Context, req *model.RequestUpdateProcessedPhoto) (error, error)
}

type facecamUseCase struct {
	// aiAdapter       adapter.AiAdapter
	// photoAdapter    adapter.PhotoAdapter
	storageAdapter  adapter.StorageAdapter
	compressAdapter adapter.CompressAdapter
	uploadProducer  producer.UploadProducer
	logs            logger.Log
}

func NewFacecamUseCase(
	// photoAdapter adapter.PhotoAdapter,
	// aiAdapter adapter.AiAdapter,
	storageAdapter adapter.StorageAdapter, compressAdapter adapter.CompressAdapter,
	uploadProducer producer.UploadProducer, logs logger.Log) FacecamUseCase {
	return &facecamUseCase{
		// photoAdapter: photoAdapter,
		// aiAdapter:       aiAdapter,
		storageAdapter:  storageAdapter,
		compressAdapter: compressAdapter,
		uploadProducer:  uploadProducer,
		logs:            logs,
	}
}

func (u *facecamUseCase) UploadFacecam(ctx context.Context, file *multipart.FileHeader, userId, creatorId string) error {
	start := time.Now()

	srcFile, err := file.Open()
	if err != nil {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Cannot open uploaded file")
	}
	defer srcFile.Close()

	// Ambil buffer 512KB dari pool
	peekBuf := peekBufPool.Get().([]byte)
	defer peekBufPool.Put(peekBuf)

	// Baca sebagian file untuk seed + streaming awal
	n, err := io.ReadFull(srcFile, peekBuf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return helper.NewUseCaseWithInternalError(errorcode.ErrInvalidArgument, "Failed to peek image data", err)
	}

	// Streaming hash & upload
	hasher := sha256.New()
	stream := io.TeeReader(io.MultiReader(bytes.NewReader(peekBuf[:n]), srcFile), hasher)

	facecamID := ulid.Make().String()
	uploadPath := fmt.Sprintf("facecam/%s/original", facecamID)
	uploaded, err := u.storageAdapter.UploadOriginalFileWithoutMultipart(ctx, file, io.NopCloser(stream), uploadPath, facecamID)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "Failed to upload file", err)
	}

	checksum := fmt.Sprintf("%x", hasher.Sum(nil))
	now := time.Now()

	log.Printf("Ini adalah uploaded file name %s", uploaded.Filename)
	newFacecam := &entity.Facecam{
		Id:         facecamID,
		UserId:     userId,
		FileName:   uploaded.Filename,
		FileKey:    uploaded.FileKey,
		Title:      uploaded.Filename,
		Size:       uploaded.Size,
		Checksum:   checksum,
		Url:        uploaded.URL,
		OriginalAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// if err := u.photoAdapter.CreateFacecam(ctx, newFacecam); err != nil {
	// 	return helper.WrapInternalServerError(u.logs, "Failed to save facecam metadata", err)
	// }

	if err := u.uploadProducer.CreateFacecam(ctx, newFacecam); err != nil {
		return helper.WrapInternalServerError(u.logs, "Failed to save facecam metadata", err)
	}

	u.logs.Log(fmt.Sprintf("✅ Facecam uploaded in %v: %s", time.Since(start), uploaded.URL))

	request := &model.ProcessFacecam{
		UserId:    userId,
		CreatorId: creatorId,
		FileURL:   uploaded.URL,
	}

	// go u.aiAdapter.ProcessFacecam(ctx, request)
	if err := u.uploadProducer.ProcessFacecam(ctx, request); err != nil {
		u.logs.CustomError("failed to process facecam", err)
	}

	return nil
}
