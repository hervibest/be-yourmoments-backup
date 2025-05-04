package converter

import (
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/nullable"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
)

func PhotoWithDetailsToGRPC(photoWithDetails *[]*entity.PhotoWithDetail, generateCDN func(string) string) *photopb.GetPhotoWithDetailsResponse {
	photoWithDetailsPBs := make([]*photopb.Photo, len(*photoWithDetails))
	for _, photoWithDetail := range *photoWithDetails {
		photoDetailPB := &photopb.PhotoDetail{
			FileName: photoWithDetail.FileName,
			// FileKey:         photoWithDetail.FileKey,
			Size:            photoWithDetail.Size,
			Type:            photoWithDetail.Type,
			Width:           photoWithDetail.Width,
			Height:          photoWithDetail.Height,
			YourMomentsType: string(photoWithDetail.YourMomentsType),
		}
		photoWithDetailsPB := &photopb.Photo{
			Id:        photoWithDetail.Id,
			CreatorId: photoWithDetail.CreatorId,
			Title:     photoWithDetail.Title,
			Price:     photoWithDetail.Price,
			PriceStr:  photoWithDetail.PriceStr,
			Latitude:  nullable.SQLToProtoDouble(photoWithDetail.Latitude),
			Url:       generateCDN(photoWithDetail.FileKey),

			Longitude:   nullable.SQLToProtoDouble(photoWithDetail.Longitude),
			Description: nullable.SQLToProtoString(photoWithDetail.Description),
			OriginalAt: &timestamppb.Timestamp{
				Seconds: photoWithDetail.OriginalAt.Unix(),
				Nanos:   int32(photoWithDetail.OriginalAt.Nanosecond()),
			},
			CreatedAt: &timestamppb.Timestamp{
				Seconds: photoWithDetail.CreatedAt.Unix(),
				Nanos:   int32(photoWithDetail.CreatedAt.Nanosecond()),
			},
			UpdatedAt: &timestamppb.Timestamp{
				Seconds: photoWithDetail.UpdatedAt.Unix(),
				Nanos:   int32(photoWithDetail.UpdatedAt.Nanosecond()),
			},

			Detail: photoDetailPB,
		}
		photoWithDetailsPBs = append(photoWithDetailsPBs, photoWithDetailsPB)
	}

	response := &photopb.GetPhotoWithDetailsResponse{
		Status:           int64(codes.OK),
		PhotoWithDetails: photoWithDetailsPBs,
	}

	return response
}
