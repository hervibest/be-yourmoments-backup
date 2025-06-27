package usecase

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/enum"
	errorcode "github.com/hervibest/be-yourmoments-backup/upload-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/nullable"
	producer "github.com/hervibest/be-yourmoments-backup/upload-svc/internal/messaging"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/model"

	_ "image/jpeg"
	_ "image/png"

	"github.com/oklog/ulid/v2"
)

type PhotoUsecase interface {
	UploadPhoto(ctx context.Context, file *multipart.FileHeader, request *model.CreatePhotoRequest) error
	BulkUploadPhoto(ctx context.Context, files []*multipart.FileHeader, request *model.CreatePhotoRequest) error
}

type photoUsecase struct {
	photoAdapter adapter.PhotoAdapter
	// aiAdapter       adapter.AiAdapter
	storageAdapter  adapter.StorageAdapter
	compressAdapter adapter.CompressAdapter
	uploadProducer  producer.UploadProducer
	logs            logger.Log
}

func NewPhotoUsecase(photoAdapter adapter.PhotoAdapter,
	// aiAdapter adapter.AiAdapter,
	storageAdapter adapter.StorageAdapter, compressAdapter adapter.CompressAdapter,
	uploadProducer producer.UploadProducer, logs logger.Log) PhotoUsecase {
	return &photoUsecase{
		photoAdapter: photoAdapter,
		// aiAdapter:       aiAdapter,
		storageAdapter:  storageAdapter,
		compressAdapter: compressAdapter,
		uploadProducer:  uploadProducer,
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
	startTotal := time.Now()

	if request.Latitude != nil && request.Longitude == nil {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Longitude is required")
	} else if request.Latitude == nil && request.Longitude != nil {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Latitude is required")
	}

	startOpen := time.Now()
	srcFile, err := file.Open()
	if err != nil {
		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Cannot open uploaded file")
	}
	defer srcFile.Close()
	u.logs.Log(fmt.Sprintf("⏱️ File open took: %v", time.Since(startOpen)))

	peekBuf := peekBufPool.Get().([]byte)
	defer peekBufPool.Put(peekBuf)

	startPeek := time.Now()
	n, err := io.ReadFull(srcFile, peekBuf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return helper.NewUseCaseWithInternalError(errorcode.ErrInvalidArgument, "Cannot read image header", err)
	}
	u.logs.Log(fmt.Sprintf("⏱️ Peek + buffer read took: %v", time.Since(startPeek)))

	startDecode := time.Now()
	imgConfig, format, err := image.DecodeConfig(bytes.NewReader(peekBuf[:n]))
	if err != nil {
		return helper.NewUseCaseWithInternalError(errorcode.ErrInvalidArgument, "Not a valid image", err)
	}
	u.logs.Log(fmt.Sprintf("⏱️ Image decode took: %v", time.Since(startDecode)))

	imageType := strings.ToUpper(format)
	if imageType == "JPEG" {
		imageType = "JPG"
	}

	hasher := sha256.New()
	multiReader := io.MultiReader(bytes.NewReader(peekBuf[:n]), srcFile)
	teeReader := io.TeeReader(multiReader, hasher)

	startUpload := time.Now()

	photoID := ulid.Make().String()
	uploadPath := fmt.Sprintf("photo/%s/original", photoID)
	upload, err := u.storageAdapter.UploadOriginalFileWithoutMultipart(ctx, file, io.NopCloser(teeReader), uploadPath, photoID)
	if err != nil {
		return helper.WrapInternalServerError(u.logs, "error when uploading file", err)
	}
	u.logs.Log(fmt.Sprintf("⏱️ Upload (original) took: %v", time.Since(startUpload)))

	checksum := fmt.Sprintf("%x", hasher.Sum(nil))
	now := time.Now()

	// Buat entitas Photo baru
	newPhoto := &entity.Photo{
		Id:            photoID,
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

	u.logs.CustomLog("Decoded image format:", format)
	u.logs.Log(fmt.Sprintf("Decoded image resolution: %d * %d", imgConfig.Width, imgConfig.Height))

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

	startDB := time.Now()
	if err := u.photoAdapter.CreatePhoto(ctx, newPhoto, newPhotoDetail); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to create photo :", err)
	}
	u.logs.Log(fmt.Sprintf("⏱️ Insert to DB took: %v", time.Since(startDB)))

	u.logs.Log(fmt.Sprintf("⏱️ TOTAL sync flow took: %v", time.Since(startTotal)))

	//
	// Goroutine: Async compression
	go func() {
		compressStart := time.Now()

		reopened, err := file.Open()
		if err != nil {
			u.logs.CustomError("failed to reopen for compression: %v", err)
			return
		}
		defer reopened.Close()

		filename, tmpFilePath, err := u.compressAdapter.CompressImageToTempFile(file.Filename, reopened)
		if err != nil {
			u.logs.CustomError("failed to compress image: %v", err)
			return
		}
		defer os.Remove(tmpFilePath)

		fileComp, err := os.Open(tmpFilePath)
		if err != nil {
			u.logs.CustomError("failed to open compressed file: %v", err)
			return
		}
		defer fileComp.Close()

		stat, _ := fileComp.Stat()
		header := &multipart.FileHeader{
			Filename: filename, // ✔️ gunakan nama file ULID.jpg
			Header:   textproto.MIMEHeader{},
			Size:     stat.Size(),
		}
		header.Header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))

		compressedPath := fmt.Sprintf("photo/%s/compressed", newPhoto.Id)
		compressedPhoto, err := u.storageAdapter.UploadCompressedFile(ctx, header, fileComp, compressedPath)
		if err != nil {
			u.logs.CustomError("failed to upload compressed file: %v", err)
			return
		}

		log.Printf("ini adalah filename dan file key %s : %s", compressedPhoto.Filename, compressedPhoto.FileKey)

		compressedPhotoDetail := &entity.PhotoDetail{
			Id:              ulid.Make().String(),
			PhotoId:         newPhoto.Id,
			FileName:        upload.Filename,
			FileKey:         compressedPhoto.FileKey,
			Size:            compressedPhoto.Size,
			Url:             compressedPhoto.URL,
			YourMomentsType: enum.YourMomentTypeCompressed,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		compDB := time.Now()
		if err := u.photoAdapter.UpdatePhotoDetail(ctx, compressedPhotoDetail); err != nil {
			u.logs.CustomError("failed to update compressed photo detail: %v", err)
			return
		}
		u.logs.Log(fmt.Sprintf("⏱️ Insert compressed photo detail took: %v", time.Since(compDB)))

		u.logs.Log(fmt.Sprintf("⏱️ TOTAL compression flow took: %v", time.Since(compressStart)))

		request := &model.ProcessPhoto{
			PhotoId:          newPhoto.Id,
			CreatorId:        newPhoto.CreatorId,
			FileURL:          compressedPhoto.URL,
			OriginalFilename: compressedPhoto.Filename,
		}

		// u.aiAdapter.ProcessPhoto(ctx, request)
		if err := u.uploadProducer.ProcessPhoto(ctx, request); err != nil {
			u.logs.CustomError("failed to produce process photo :", err)
		}
	}()

	return nil
}

var peekBufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 512*1024) // 512 KB buffer
	},
}

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

	photoEntities := make([]*entity.Photo, 0, len(files))
	photoDetailEntities := make([]*entity.PhotoDetail, 0, len(files))

	for _, file := range files {
		srcFile, err := file.Open()
		if err != nil {
			return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Cannot open uploaded file")
		}
		defer srcFile.Close()

		// Gunakan peek buffer dari sync.Pool
		peekBuf := peekBufPool.Get().([]byte)
		defer peekBufPool.Put(peekBuf)

		n, err := io.ReadFull(srcFile, peekBuf)
		if err != nil && err != io.ErrUnexpectedEOF {
			return helper.NewUseCaseWithInternalError(errorcode.ErrInvalidArgument, "Cannot read image header", err)
		}

		// Decode metadata image dari awal buffer
		imgConfig, format, err := image.DecodeConfig(bytes.NewReader(peekBuf[:n]))
		if err != nil {
			return helper.NewUseCaseWithInternalError(errorcode.ErrInvalidArgument, "Not a valid image", err)
		}

		imageType := strings.ToUpper(format)
		if imageType == "JPEG" {
			imageType = "JPG"
		}

		// Buat checksum saat streaming
		hasher := sha256.New()
		multiReader := io.MultiReader(
			bytes.NewReader(peekBuf[:n]),
			srcFile,
		)
		teeReader := io.TeeReader(multiReader, hasher)

		photoID := ulid.Make().String()
		// Upload langsung ke storage
		uploadPath := fmt.Sprintf("photo/%s/original", photoID)
		upload, err := u.storageAdapter.UploadOriginalFileWithoutMultipart(ctx, file, io.NopCloser(teeReader), uploadPath, photoID)
		if err != nil {
			return helper.WrapInternalServerError(u.logs, "error when uploading file", err)
		}

		checksum := fmt.Sprintf("%x", hasher.Sum(nil))
		now := time.Now()

		newPhoto := &entity.Photo{
			Id:            photoID,
			UserId:        request.UserId,
			CreatorId:     request.CreatorId,
			BulkPhotoId:   nullable.ToSQLString(&bulkPhoto.Id),
			Title:         upload.Filename,
			CollectionUrl: upload.URL,
			Price:         request.Price,
			PriceStr:      request.PriceStr,
			OriginalAt:    now,
			CreatedAt:     now,
			UpdatedAt:     now,
			Latitude:      nullable.ToSQLFloat64(request.Latitude),
			Longitude:     nullable.ToSQLFloat64(request.Longitude),
			Description:   nullable.ToSQLString(request.Description),
		}
		photoEntities = append(photoEntities, newPhoto)

		newPhotoDetail := &entity.PhotoDetail{
			Id:              ulid.Make().String(),
			PhotoId:         newPhoto.Id,
			FileName:        upload.Filename,
			FileKey:         upload.FileKey,
			Size:            upload.Size,
			Type:            imageType,
			Checksum:        checksum,
			Width:           imgConfig.Width,
			Height:          imgConfig.Height,
			Url:             upload.URL,
			YourMomentsType: enum.YourMomentTypeCollection,
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		photoDetailEntities = append(photoDetailEntities, newPhotoDetail)

		u.logs.Log(fmt.Sprintf("Uploaded photo %s (%dx%d) [%s]", upload.Filename, imgConfig.Width, imgConfig.Height, imageType))
	}

	if err := u.photoAdapter.CreatePhotos(ctx, bulkPhoto, &photoEntities, &photoDetailEntities); err != nil {
		return helper.WrapInternalServerError(u.logs, "failed to create photo entities", err)
	}

	// go func() {
	// 	if err := u.aiAdapter.ProcessBulkPhoto(ctx, bulkPhoto, &photoEntities); err != nil {
	// 		log.Printf("failed to process bulk photo via grpc: %v", err)
	// 	}
	// }()

	if err := u.uploadProducer.ProcessBulkPhoto(ctx, bulkPhoto, &photoEntities); err != nil {
		u.logs.CustomError("failed to produce process bulk photo", err)
	}

	return nil
}
