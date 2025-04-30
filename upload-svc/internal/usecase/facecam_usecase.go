package usecase

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"time"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/entity"
	errorcode "github.com/hervibest/be-yourmoments-backup/upload-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/logger"

	"github.com/oklog/ulid/v2"
)

type FacecamUseCase interface {
	UploadFacecam(ctx context.Context, file *multipart.FileHeader, userId string) error
	// UpdateProcessedPhoto(ctx context.Context, req *model.RequestUpdateProcessedPhoto) (error, error)
}

type facecamUseCase struct {
	aiAdapter       adapter.AiAdapter
	photoAdapter    adapter.PhotoAdapter
	storageAdapter  adapter.StorageAdapter
	compressAdapter adapter.CompressAdapter
	logs            logger.Log
}

func NewFacecamUseCase(aiAdapter adapter.AiAdapter, photoAdapter adapter.PhotoAdapter,
	storageAdapter adapter.StorageAdapter, compressAdapter adapter.CompressAdapter,
	logs logger.Log) FacecamUseCase {
	return &facecamUseCase{
		aiAdapter:       aiAdapter,
		photoAdapter:    photoAdapter,
		storageAdapter:  storageAdapter,
		compressAdapter: compressAdapter,
		logs:            logs,
	}
}

func (u *facecamUseCase) UploadFacecam(ctx context.Context, file *multipart.FileHeader, userId string) error {
	uploadFile, err := file.Open()
	if err != nil {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Cannot process uploaded file")
	}

	data, err := io.ReadAll(uploadFile)
	if err != nil {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Cannot read uploaded file")
	}

	uploadFile.Close()

	readerForUpload := bytes.NewReader(data)
	wrappedReader := nopReadSeekCloser{readerForUpload}

	checksum := fmt.Sprintf("%x", sha256.Sum256(data))

	go func() {
		_, filePath, err := u.compressAdapter.CompressImage(file, wrappedReader, "facecam")
		if err != nil {
			u.logs.CustomError("Error compressing images: %v", err)
			return
		}

		fileComp, err := os.Open(filePath)
		if err != nil {
			u.logs.CustomError("Error opening file: %v", err)
			return
		}
		defer fileComp.Close()

		fileInfo, err := fileComp.Stat()
		if err != nil {
			u.logs.CustomError("Error stating file: %v", err)
			return
		}

		mimeHeader := make(textproto.MIMEHeader)
		mimeHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileInfo.Name()))

		fileHeader := &multipart.FileHeader{
			Filename: fileInfo.Name(),
			Header:   mimeHeader,
			Size:     fileInfo.Size(),
		}

		uploadPath := "facecam/compressed"
		compressedPhoto, err := u.storageAdapter.UploadFile(ctx, fileHeader, fileComp, uploadPath)
		if err != nil {
			u.logs.CustomError("Error uploading file: %v", err)
			return
		}

		newFacecam := &entity.Facecam{
			Id:         ulid.Make().String(),
			UserId:     userId,
			FileName:   compressedPhoto.Filename,
			FileKey:    compressedPhoto.FileKey,
			Title:      compressedPhoto.Filename,
			Size:       compressedPhoto.Size,
			Checksum:   checksum,
			Url:        compressedPhoto.URL,
			OriginalAt: time.Now(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := u.photoAdapter.CreateFacecam(ctx, newFacecam); err != nil {
			u.logs.CustomError("Error creating facecam: %v", err)
			return
		}

		if err := os.Remove(filePath); err != nil {
			u.logs.CustomError("Gagal menghapus file: %v", err)
		} else {
			u.logs.CustomLog("File sementara berhasil dihapus: %s", filePath)
		}

		u.aiAdapter.ProcessFacecam(ctx, userId, compressedPhoto.URL)
	}()

	return nil

}
