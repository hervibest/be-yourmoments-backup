package usecase

import (
	"be-yourmoments/upload-svc/internal/adapter"
	"be-yourmoments/upload-svc/internal/entity"
	"be-yourmoments/upload-svc/internal/enum"
	errorcode "be-yourmoments/upload-svc/internal/enum/error"
	"be-yourmoments/upload-svc/internal/helper"
	"be-yourmoments/upload-svc/internal/helper/logger"
	"be-yourmoments/upload-svc/internal/helper/nullable"
	"be-yourmoments/upload-svc/internal/model"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/oklog/ulid/v2"
)

type PhotoUsecase interface {
	UploadPhoto(ctx context.Context, file *multipart.FileHeader, request *model.CreatePhotoRequest) error
	BulkUploadPhoto(ctx context.Context, files []*multipart.FileHeader, request *model.CreatePhotoRequest) error
}

type photoUsecase struct {
	aiAdapter       adapter.AiAdapter
	photoAdapter    adapter.PhotoAdapter
	storageAdapter  adapter.StorageAdapter
	compressAdapter adapter.CompressAdapter
	logs            *logger.Log
}

func NewPhotoUsecase(aiAdapter adapter.AiAdapter, photoAdapter adapter.PhotoAdapter,
	storageAdapter adapter.StorageAdapter, compressAdapter adapter.CompressAdapter,
	logs *logger.Log) PhotoUsecase {
	return &photoUsecase{
		aiAdapter:       aiAdapter,
		photoAdapter:    photoAdapter,
		storageAdapter:  storageAdapter,
		compressAdapter: compressAdapter,
		logs:            logs,
	}
}

type nopReadSeekCloser struct {
	*bytes.Reader
}

func (n nopReadSeekCloser) Close() error {
	return nil
}

func (u *photoUsecase) UploadPhoto(ctx context.Context, file *multipart.FileHeader, request *model.CreatePhotoRequest) error {
	if request.Latitude != nil && request.Longitude == nil {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Longitude is required")
	} else if request.Latitude == nil && request.Longitude != nil {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Latitude is required")
	}

	uploadFile, err := file.Open()
	if err != nil {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Cannot process uploaded file")
	}

	data, err := io.ReadAll(uploadFile)
	if err != nil {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Cannot read uploaded file")
	}

	defer uploadFile.Close()

	readerForUpload := bytes.NewReader(data)
	wrappedReader := nopReadSeekCloser{readerForUpload}

	upload, err := u.storageAdapter.UploadFile(ctx, file, wrappedReader, "photo")
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "error when uploading file :", err)
	}

	// Buat entitas Photo baru
	newPhoto := &entity.Photo{
		Id:            ulid.Make().String(),
		UserId:        request.UserId,
		CreatorId:     request.CreatorId,
		Title:         upload.Filename,
		CollectionUrl: upload.URL,
		Price:         request.Price,
		PriceStr:      request.PriceStr,
		OriginalAt:    time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Latitude:      nullable.ToSQLFloat64(request.Latitude),
		Longitude:     nullable.ToSQLFloat64(request.Longitude),
		Description:   nullable.ToSQLString(request.Description),
	}

	// Lanjutkan ke database insert, dll

	readerForDecode := bytes.NewReader(data)
	imgConfig, format, err := image.DecodeConfig(readerForDecode)
	if err != nil {
		return helper.NewUseCaseWithInternalError(errorcode.ErrInvalidArgument, "Not a valid image", err)
	}

	u.logs.CustomLog("Decoded image format:", format)
	u.logs.Log(fmt.Sprintf("Decoded image resolution: %d * %d", imgConfig.Width, imgConfig.Height))

	var imageType string
	if format == "jpeg" {
		imageType = "JPG"
	} else {
		imageType = strings.ToUpper(format)
	}

	checksum := fmt.Sprintf("%x", sha256.Sum256(data))

	newPhotoDetail := &entity.PhotoDetail{
		Id:              ulid.Make().String(),
		PhotoId:         newPhoto.Id,
		FileName:        upload.Filename,
		FileKey:         upload.FileKey,
		Size:            upload.Size,
		Type:            imageType,
		Checksum:        checksum,
		Width:           imgConfig.Width,  // disesuaikan tipe data jika perlu
		Height:          imgConfig.Height, // disesuaikan tipe data jika perlu
		Url:             upload.URL,
		YourMomentsType: enum.YourMomentTypeCollection,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := u.photoAdapter.CreatePhoto(ctx, newPhoto, newPhotoDetail); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to create photo :", err)
	}

	/* TO DO
	1. Pastikan I/O untuk open file dilakukan secara efisien disini
	2. Pastikan persitensy hasil kompres dilakukan dengan baik

	*/
	go func() {
		// Should be queued ? using go routine ? Bisa dipikirkan nanti
		_, filePath, err := u.compressAdapter.CompressImage(file, wrappedReader, "photo")
		if err != nil {
			u.logs.CustomError("failed to compress image: %v", err)
			return
		}

		fileComp, err := os.Open(filePath)
		if err != nil {
			u.logs.CustomError("failed to open compressed file: %v", err)
			return
		}
		defer fileComp.Close()

		fileInfo, err := fileComp.Stat()
		if err != nil {
			u.logs.CustomError("failed to stating compressed file: %v", err)
			return
		}

		mimeHeader := make(textproto.MIMEHeader)
		mimeHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileInfo.Name()))

		fileHeader := &multipart.FileHeader{
			Filename: fileInfo.Name(),
			Header:   mimeHeader,
			Size:     fileInfo.Size(),
		}

		uploadPath := "photo/compressed"
		compressedPhoto, err := u.storageAdapter.UploadFile(ctx, fileHeader, fileComp, uploadPath)
		if err != nil {
			u.logs.CustomError("failed to upload compressed file: %v", err)
			return
		}

		compressedPhotoDetail := &entity.PhotoDetail{
			Id:              ulid.Make().String(),
			PhotoId:         newPhoto.Id,
			FileName:        compressedPhoto.Filename,
			FileKey:         compressedPhoto.FileKey,
			Size:            compressedPhoto.Size,
			Url:             compressedPhoto.URL,
			YourMomentsType: enum.YourMomentTypeCompressed,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if err := u.photoAdapter.UpdatePhotoDetail(ctx, compressedPhotoDetail); err != nil {
			u.logs.CustomError("failed to update compressed photo detail: %v", err)
			return
		}

		if err := os.Remove(filePath); err != nil {
			u.logs.CustomError("failed to remove file: %v", err)
		} else {
			u.logs.CustomLog("remove file sucess : %s", filePath)
		}

		u.aiAdapter.ProcessPhoto(ctx, newPhoto.Id, compressedPhoto.URL)
	}()

	return nil
}

// ISSUE #2 Dont forget to add bulk_photo_id for each photo entity
func (u *photoUsecase) BulkUploadPhoto(ctx context.Context, files []*multipart.FileHeader, request *model.CreatePhotoRequest) error {
	if request.Latitude != nil && request.Longitude == nil {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Longitude is required")
	} else if request.Latitude == nil && request.Longitude != nil {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Latitude is required")
	}

	bulkPhoto := &entity.BulkPhoto{
		Id:              ulid.Make().String(),
		CreatorId:       request.CreatorId,
		BulkPhotoStatus: enum.BulkPhotoStatusProcessed,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	photoEntities := make([]*entity.Photo, 0)
	photoDetailEntities := make([]*entity.PhotoDetail, 0)

	// Initiate
	for _, file := range files {
		uploadFile, err := file.Open()
		if err != nil {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Cannot process uploaded file")
		}

		data, err := io.ReadAll(uploadFile)
		if err != nil {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Cannot read uploaded file")
		}

		defer uploadFile.Close()

		readerForUpload := bytes.NewReader(data)
		wrappedReader := nopReadSeekCloser{readerForUpload}

		upload, err := u.storageAdapter.UploadFile(ctx, file, wrappedReader, "photo")
		if err != nil {
			return helper.WrapInternalServerError(u.logs, "error when uploading file :", err)
		}

		// Buat entitas Photo baru
		newPhoto := &entity.Photo{
			Id:            ulid.Make().String(),
			UserId:        request.UserId,
			CreatorId:     request.CreatorId,
			BulkPhotoId:   nullable.ToSQLString(&bulkPhoto.Id),
			Title:         upload.Filename,
			CollectionUrl: upload.URL,
			Price:         request.Price,
			PriceStr:      request.PriceStr,
			OriginalAt:    time.Now(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Latitude:      nullable.ToSQLFloat64(request.Latitude),
			Longitude:     nullable.ToSQLFloat64(request.Longitude),
			Description:   nullable.ToSQLString(request.Description),
		}

		photoEntities = append(photoEntities, newPhoto)
		readerForDecode := bytes.NewReader(data)
		imgConfig, format, err := image.DecodeConfig(readerForDecode)
		if err != nil {
			return helper.NewUseCaseWithInternalError(errorcode.ErrInvalidArgument, "Not a valid image", err)
		}

		u.logs.CustomLog("Decoded image format:", format)
		u.logs.Log(fmt.Sprintf("Decoded image resolution: %d * %d", imgConfig.Width, imgConfig.Height))

		var imageType string
		if format == "jpeg" {
			imageType = "JPG"
		} else {
			imageType = strings.ToUpper(format)
		}

		checksum := fmt.Sprintf("%x", sha256.Sum256(data))

		newPhotoDetail := &entity.PhotoDetail{
			Id:              ulid.Make().String(),
			PhotoId:         newPhoto.Id,
			FileName:        upload.Filename,
			FileKey:         upload.FileKey,
			Size:            upload.Size,
			Type:            imageType,
			Checksum:        checksum,
			Width:           imgConfig.Width,  // disesuaikan tipe data jika perlu
			Height:          imgConfig.Height, // disesuaikan tipe data jika perlu
			Url:             upload.URL,
			YourMomentsType: enum.YourMomentTypeCollection,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		photoDetailEntities = append(photoDetailEntities, newPhotoDetail)
	}

	if err := u.photoAdapter.CreatePhotos(ctx, bulkPhoto, &photoEntities, &photoDetailEntities); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to create photo :", err)
	}

	return nil
}
