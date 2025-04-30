package usecase_test

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"mime/multipart"
	"testing"

	"github.com/golang/mock/gomock"
	mocksadapter "github.com/hervibest/be-yourmoments-backup/upload-svc/internal/mocks/adapter"
	mockslogger "github.com/hervibest/be-yourmoments-backup/upload-svc/internal/mocks/logger"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/usecase"
	"github.com/stretchr/testify/require"
)

func TestUploadPhoto_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocksadapter.NewMockStorageAdapter(ctrl)
	mockPhoto := mocksadapter.NewMockPhotoAdapter(ctrl)
	mockCompress := mocksadapter.NewMockCompressAdapter(ctrl)
	mockAI := mocksadapter.NewMockAiAdapter(ctrl)
	mockLogger := mockslogger.NewMockLog(ctrl)

	usecase := usecase.NewPhotoUsecase(mockAI, mockPhoto, mockStorage, mockCompress, mockLogger)
	// Simulasi file jpeg kecil
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, nil)
	fileData := buf.Bytes()

	fileReader := multipart.FileHeader{
		Filename: "test.jpg",
		Size:     int64(len(fileData)),
	}

	request := &model.CreatePhotoRequest{
		UserId:    "user123",
		CreatorId: "creator456",
		Latitude:  floatPtr(7.1),
		Longitude: floatPtr(110.2),
	}

	// Expect upload
	mockStorage.EXPECT().
		UploadFile(gomock.Any(), &fileReader, gomock.Any(), "photo").
		Return(&model.MinioFileResponse{
			Filename: "test.jpg",
			URL:      "http://example.com/photo.jpg",
			FileKey:  "photo/test.jpg",
			Size:     int64(len(fileData)),
		}, nil)

	mockPhoto.EXPECT().
		CreatePhoto(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	// Optional compress, bisa dilewatkan atau disimulasikan

	err := usecase.UploadPhoto(context.Background(), &fileReader, request)
	require.NoError(t, err)
}

func floatPtr(f float64) *float64 {
	return &f
}
