package adapter

import (
	"be-yourmoments/upload-svc/internal/entity"
	"be-yourmoments/upload-svc/internal/helper"
	"be-yourmoments/upload-svc/internal/helper/discovery"
	"be-yourmoments/upload-svc/internal/helper/logger"
	"be-yourmoments/upload-svc/internal/helper/nullable"
	"context"

	photopb "github.com/be-yourmoments/pb/photo"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PhotoAdapter interface {
	CreatePhoto(ctx context.Context, photo *entity.Photo, facecam *entity.PhotoDetail) error
	CreatePhotos(ctx context.Context, bulkPhoto *entity.BulkPhoto, photos *[]*entity.Photo, photoDetails *[]*entity.PhotoDetail) error
	UpdatePhotoDetail(ctx context.Context, facecam *entity.PhotoDetail) error
	CreateFacecam(ctx context.Context, facecam *entity.Facecam) error
}

type photoAdapter struct {
	client photopb.PhotoServiceClient
	logs   *logger.Log
}

func NewPhotoAdapter(ctx context.Context, registry discovery.Registry, logs *logger.Log) (PhotoAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "photo-svc-grpc", registry)
	if err != nil {
		return nil, err
	}
	logs.Log("successfuly connected to photo-svc-grpc")
	client := photopb.NewPhotoServiceClient(conn)

	return &photoAdapter{
		client: client,
		logs:   logs,
	}, nil
}

func (a *photoAdapter) CreatePhoto(ctx context.Context, photo *entity.Photo, facecam *entity.PhotoDetail) error {
	photoDetailPb := &photopb.PhotoDetail{
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

	photoPb := &photopb.Photo{
		Id:            photo.Id,
		UserId:        photo.UserId,
		CreatorId:     photo.CreatorId,
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

		Detail:      photoDetailPb,
		Latitude:    nullable.SQLToProtoDouble(photo.Latitude),
		Longitude:   nullable.SQLToProtoDouble(photo.Longitude),
		Description: nullable.SQLToProtoString(photo.Description),
	}

	pbRequest := &photopb.CreatePhotoRequest{
		Photo: photoPb,
	}

	_, err := a.client.CreatePhoto(context.Background(), pbRequest)
	if err != nil {
		return helper.FromGRPCError(err)
	}

	return nil
}

func (a *photoAdapter) UpdatePhotoDetail(ctx context.Context, facecam *entity.PhotoDetail) error {
	photoDetailPb := &photopb.PhotoDetail{
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
	pbRequest := &photopb.UpdatePhotoDetailRequest{
		PhotoDetail: photoDetailPb,
	}

	_, err := a.client.UpdatePhotoDetail(context.Background(), pbRequest)
	if err != nil {
		return helper.FromGRPCError(err)
	}

	return nil
}

func (a *photoAdapter) CreateFacecam(ctx context.Context, facecam *entity.Facecam) error {
	facecamPb := &photopb.Facecam{
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

	pbRequest := &photopb.CreateFacecamRequest{
		Facecam: facecamPb,
	}

	_, err := a.client.CreateFacecam(context.Background(), pbRequest)
	if err != nil {
		return helper.FromGRPCError(err)
	}

	return nil
}

func (a *photoAdapter) CreatePhotos(ctx context.Context, bulkPhoto *entity.BulkPhoto, photos *[]*entity.Photo, photoDetails *[]*entity.PhotoDetail) error {
	photoPbs := make([]*photopb.Photo, 0)

	photoBukPb := &photopb.BulkPhoto{
		Id:              bulkPhoto.Id,
		CreatorId:       bulkPhoto.CreatorId,
		BulkPhotoStatus: string(bulkPhoto.BulkPhotoStatus),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: bulkPhoto.CreatedAt.Unix(),
			Nanos:   int32(bulkPhoto.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: bulkPhoto.UpdatedAt.Unix(),
			Nanos:   int32(bulkPhoto.UpdatedAt.Nanosecond()),
		},
	}

	for idx, photoDetail := range *photoDetails {
		photoDetailPb := &photopb.PhotoDetail{
			Id:              photoDetail.Id,
			PhotoId:         photoDetail.PhotoId,
			FileName:        photoDetail.FileName,
			FileKey:         photoDetail.FileKey,
			Size:            photoDetail.Size,
			Type:            photoDetail.Type,
			Checksum:        photoDetail.Checksum,
			Width:           int32(photoDetail.Width),
			Height:          int32(photoDetail.Height),
			Url:             photoDetail.Url,
			YourMomentsType: string(photoDetail.YourMomentsType),
			CreatedAt: &timestamppb.Timestamp{
				Seconds: photoDetail.CreatedAt.Unix(),
				Nanos:   int32(photoDetail.CreatedAt.Nanosecond()),
			},
			UpdatedAt: &timestamppb.Timestamp{
				Seconds: photoDetail.UpdatedAt.Unix(),
				Nanos:   int32(photoDetail.UpdatedAt.Nanosecond()),
			},
		}

		photoPb := &photopb.Photo{
			Id:            (*photos)[idx].Id,
			UserId:        (*photos)[idx].UserId,
			CreatorId:     (*photos)[idx].CreatorId,
			Title:         (*photos)[idx].Title,
			BulkPhotoId:   nullable.SQLToProtoString((*photos)[idx].BulkPhotoId),
			CollectionUrl: (*photos)[idx].CollectionUrl,
			Price:         int32((*photos)[idx].Price),
			PriceStr:      (*photos)[idx].PriceStr,

			OriginalAt: &timestamppb.Timestamp{
				Seconds: (*photos)[idx].OriginalAt.Unix(),
				Nanos:   int32((*photos)[idx].OriginalAt.Nanosecond()),
			},
			CreatedAt: &timestamppb.Timestamp{
				Seconds: (*photos)[idx].CreatedAt.Unix(),
				Nanos:   int32((*photos)[idx].CreatedAt.Nanosecond()),
			},
			UpdatedAt: &timestamppb.Timestamp{
				Seconds: (*photos)[idx].UpdatedAt.Unix(),
				Nanos:   int32((*photos)[idx].UpdatedAt.Nanosecond()),
			},

			Detail:      photoDetailPb,
			Latitude:    nullable.SQLToProtoDouble((*photos)[idx].Latitude),
			Longitude:   nullable.SQLToProtoDouble((*photos)[idx].Longitude),
			Description: nullable.SQLToProtoString((*photos)[idx].Description),
		}

		photoPbs = append(photoPbs, photoPb)
	}

	pbRequest := &photopb.CreateBulkPhotoRequest{
		BulkPhoto: photoBukPb,
		Photos:    photoPbs,
	}

	_, err := a.client.CreateBulkPhoto(context.Background(), pbRequest)
	if err != nil {
		return helper.FromGRPCError(err)
	}

	return nil
}
