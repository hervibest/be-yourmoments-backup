package adapter

import (
	"be-yourmoments/upload-svc/internal/entity"
	"be-yourmoments/upload-svc/internal/helper"
	"be-yourmoments/upload-svc/internal/helper/discovery"
	"be-yourmoments/upload-svc/internal/helper/logger"
	"be-yourmoments/upload-svc/internal/helper/nullable"
	"context"

	"github.com/be-yourmoments/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PhotoAdapter interface {
	CreatePhoto(ctx context.Context, photo *entity.Photo, facecam *entity.PhotoDetail) error
	UpdatePhotoDetail(ctx context.Context, facecam *entity.PhotoDetail) error
	CreateFacecam(ctx context.Context, facecam *entity.Facecam) error
}

type photoAdapter struct {
	client pb.PhotoServiceClient
	logs   *logger.Log
}

func NewPhotoAdapter(ctx context.Context, registry discovery.Registry, logs *logger.Log) (PhotoAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "photo-svc-grpc", registry)
	if err != nil {
		return nil, err
	}
	logs.Log("successfuly connected to photo-svc-grpc")
	client := pb.NewPhotoServiceClient(conn)

	return &photoAdapter{
		client: client,
		logs:   logs,
	}, nil
}

func (a *photoAdapter) CreatePhoto(ctx context.Context, photo *entity.Photo, facecam *entity.PhotoDetail) error {
	facecampb := &pb.PhotoDetail{
		Id:              facecam.Id,
		PhotoId:         facecam.PhotoId,
		FileName:        facecam.FileName,
		FileKey:         facecam.FileKey,
		Size:            facecam.Size,
		Type:            facecam.Type,
		Checksum:        facecam.Checksum,
		Width:           int32(facecam.Width),
		Height:          int32(facecam.Height),
		Url:             facecam.Url,
		YourMomentsType: string(facecam.YourMomentsType),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: facecam.CreatedAt.Unix(),
			Nanos:   int32(facecam.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: facecam.UpdatedAt.Unix(),
			Nanos:   int32(facecam.UpdatedAt.Nanosecond()),
		},
	}

	photoPb := &pb.Photo{
		Id:            photo.Id,
		UserId:        photo.UserId,
		Title:         photo.Title,
		CollectionUrl: photo.CollectionUrl,
		Price:         int32(photo.Price),
		PriceStr:      photo.PriceStr,

		OriginalAt: &timestamppb.Timestamp{
			Seconds: photo.OriginalAt.Unix(),
			Nanos:   int32(photo.OriginalAt.Nanosecond()),
		},
		CreatedAt: &timestamppb.Timestamp{
			Seconds: photo.CreatedAt.Unix(),
			Nanos:   int32(photo.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: photo.UpdatedAt.Unix(),
			Nanos:   int32(photo.UpdatedAt.Nanosecond()),
		},

		Detail:      facecampb,
		Latitude:    nullable.SQLToProtoDouble(photo.Latitude),
		Longitude:   nullable.SQLToProtoDouble(photo.Longitude),
		Description: nullable.SQLToProtoString(photo.Description),
	}

	pbRequest := &pb.CreatePhotoRequest{
		Photo: photoPb,
	}

	_, err := a.client.CreatePhoto(context.Background(), pbRequest)
	if err != nil {
		return helper.FromGRPCError(err)
	}

	return nil
}

func (a *photoAdapter) UpdatePhotoDetail(ctx context.Context, facecam *entity.PhotoDetail) error {
	facecampb := &pb.PhotoDetail{
		Id:              facecam.Id,
		PhotoId:         facecam.PhotoId,
		FileName:        facecam.FileName,
		FileKey:         facecam.FileKey,
		Size:            facecam.Size,
		Type:            facecam.Type,
		Checksum:        facecam.Checksum,
		Width:           int32(facecam.Width),
		Height:          int32(facecam.Height),
		Url:             facecam.Url,
		YourMomentsType: string(facecam.YourMomentsType),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: facecam.CreatedAt.Unix(),
			Nanos:   int32(facecam.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: facecam.UpdatedAt.Unix(),
			Nanos:   int32(facecam.UpdatedAt.Nanosecond()),
		},
	}
	pbRequest := &pb.UpdatePhotoDetailRequest{
		PhotoDetail: facecampb,
	}

	_, err := a.client.UpdatePhotoDetail(context.Background(), pbRequest)
	if err != nil {
		return helper.FromGRPCError(err)
	}

	return nil
}

func (a *photoAdapter) CreateFacecam(ctx context.Context, facecam *entity.Facecam) error {
	facecamPb := &pb.Facecam{
		Id:       facecam.Id,
		UserId:   facecam.UserId,
		FileName: facecam.FileName,
		FileKey:  facecam.FileKey,
		Size:     facecam.Size,
		Checksum: facecam.Checksum,
		Url:      facecam.Url,
		CreatedAt: &timestamppb.Timestamp{
			Seconds: facecam.CreatedAt.Unix(),
			Nanos:   int32(facecam.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: facecam.UpdatedAt.Unix(),
			Nanos:   int32(facecam.UpdatedAt.Nanosecond()),
		},
	}

	pbRequest := &pb.CreateFacecamRequest{
		Facecam: facecamPb,
	}

	_, err := a.client.CreateFacecam(context.Background(), pbRequest)
	if err != nil {
		return helper.FromGRPCError(err)
	}

	return nil
}
