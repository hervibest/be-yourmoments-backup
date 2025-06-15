package converter

import (
	photopb "github.com/hervibest/be-yourmoments-backup/pb/photo"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/entity"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/enum"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/nullable"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/model"
	"github.com/midtrans/midtrans-go/coreapi"
)

func TransactionToResponse(transaction *entity.Transaction, redirectUrl string) *model.CreateTransactionResponse {
	return &model.CreateTransactionResponse{
		TransactionId: transaction.Id,
		SnapToken:     transaction.SnapToken.String,
		RedirectURL:   redirectUrl,
	}
}
func TransactionAndPhotoToSingleResponse(
	transactionWithDetails []*entity.TransactionWithDetail,
	photoWithDetailPBs []*photopb.Photo,
) *model.TransactionWithDetail {

	if len(transactionWithDetails) == 0 {
		return nil
	}

	// Buat index PhotoId ke Photo proto
	photoMap := make(map[string]*photopb.Photo)
	for _, photo := range photoWithDetailPBs {
		photoMap[photo.GetId()] = photo
	}

	// Gunakan data dari baris pertama sebagai data utama transaksi
	base := transactionWithDetails[0]

	result := &model.TransactionWithDetail{
		TransactionId:       base.TransactionId,
		UserId:              base.UserId,
		Status:              base.Status,
		TransactionMethodId: nullable.SQLStringToPtr(base.TransactionMethodId),
		TransactionTypeId:   nullable.SQLStringToPtr(base.TransactionTypeId),
		PaymentTypeId:       nullable.SQLStringToPtr(base.PaymentTypeId),
		PaymentAt:           base.PaymentAt,
		CheckoutAt:          base.CheckoutAt,
		Amount:              base.Amount,
		CreatedAt:           base.CreatedAt,
		UpdatedAt:           base.UpdatedAt,
		TransactionDetail:   &[]*model.TransactionDetailResponse{},
	}

	// Map untuk mengelompokkan photos berdasarkan creator_id
	detailMap := make(map[string]*model.TransactionDetailResponse)

	for _, trx := range transactionWithDetails {
		detail, exists := detailMap[trx.CreatorId]
		if !exists {
			detail = &model.TransactionDetailResponse{
				TransactionDetailId: trx.TranscationDetailId,
				CreatorId:           trx.CreatorId,
				CreatorDiscountId:   trx.CreatorDiscountId,
				IsReviewed:          trx.IsReviewed,
				Photo:               &[]*model.PhotoResponse{},
			}
			detailMap[trx.CreatorId] = detail
		}

		if photoPB, ok := photoMap[trx.PhotoId]; ok {
			*detail.Photo = append(*detail.Photo, mapPhotoPBToPhotoResponse(photoPB, trx))
		}
	}

	// Masukkan hasilnya ke slice TransactionDetail
	for _, d := range detailMap {
		*result.TransactionDetail = append(*result.TransactionDetail, d)
	}

	return result
}

// Mapping dari protobuf Photo ke response struct
func mapPhotoPBToPhotoResponse(photo *photopb.Photo, trx *entity.TransactionWithDetail) *model.PhotoResponse {
	return &model.PhotoResponse{
		PhotoId:         photo.GetId(),
		Price:           trx.Price,
		Discount:        trx.Discount,
		FinalPrice:      trx.FinalPrice,
		Url:             photo.GetUrl(),
		Title:           photo.GetTitle(),
		Latitude:        nullable.WrapDouble(photo.Latitude),
		Longitude:       nullable.WrapDouble(photo.Longitude),
		Description:     nullable.WrapString(photo.Description),
		PhotoOriginalAt: photo.GetOriginalAt().AsTime(),
		PhotoCreatedAt:  photo.GetCreatedAt().AsTime(),
		PhotoUpdatedAt:  photo.GetUpdatedAt().AsTime(),
		FileName:        photo.GetDetail().GetFileName(),
		Size:            photo.GetDetail().GetSize(),
		Type:            photo.GetDetail().GetType(),
		Width:           photo.GetDetail().GetWidth(),
		Height:          photo.GetDetail().GetHeight(),
		YourMomentsType: enum.YourMomentsType(photo.GetDetail().GetYourMomentsType()),
	}
}

func UserTransactionToResponse(userTransactions *[]*entity.Transaction) *[]*model.UserTransaction {
	userTransactionReponses := make([]*model.UserTransaction, 0)
	counter := 0
	for _, userTransaction := range *userTransactions {
		counter++
		userTransactionResponse := &model.UserTransaction{
			Id:                  userTransaction.Id,
			UserId:              userTransaction.UserId,
			Status:              userTransaction.Status,
			TransactionMethodId: nullable.SQLStringToPtr(userTransaction.TransactionMethodId),
			TransactionTypeId:   nullable.SQLStringToPtr(userTransaction.TransactionTypeId),
			PaymentTypeId:       nullable.SQLStringToPtr(userTransaction.PaymentTypeId),
			PaymentAt:           userTransaction.PaymentAt,
			CheckoutAt:          userTransaction.CheckoutAt,
			Amount:              userTransaction.Amount,
			// CreatedAt:           userTransaction.CreatedAt, 
			// UpdatedAt:           userTransaction.UpdatedAt,
		}
		userTransactionReponses = append(userTransactionReponses, userTransactionResponse)
	}
	return &userTransactionReponses
}

func WebhookReqToCheckAndUpdate(webhookReq *model.UpdateTransactionWebhookRequest) *model.CheckAndUpdateTransactionRequest {
	return &model.CheckAndUpdateTransactionRequest{
		MidtransTransactionType:   webhookReq.MidtransTransactionType,
		MidtransTransactionTime:   webhookReq.MidtransTransactionTime,
		MidtransTransactionStatus: webhookReq.MidtransTransactionStatus,
		MidtransTransactionID:     webhookReq.MidtransTransactionID,
		StatusMessage:             webhookReq.StatusMessage,
		StatusCode:                webhookReq.StatusCode,
		SignatureKey:              webhookReq.SignatureKey,
		SettlementTime:            webhookReq.SettlementTime,
		ReferenceID:               webhookReq.ReferenceID,
		PaymentType:               webhookReq.PaymentType,
		OrderID:                   webhookReq.OrderID,
		Metadata:                  webhookReq.Metadata,
		MerchantID:                webhookReq.MerchantID,
		GrossAmount:               webhookReq.GrossAmount,
		FraudStatus:               webhookReq.FraudStatus,
		ExpiryTime:                webhookReq.ExpiryTime,
		Currency:                  webhookReq.Currency,
		Acquirer:                  webhookReq.Acquirer,
		Body:                      webhookReq.Body,
	}
}

func SchedulerReqToCheckAndUpdate(schedulerReq *coreapi.TransactionStatusResponse, body []byte) *model.CheckAndUpdateTransactionRequest {
	return &model.CheckAndUpdateTransactionRequest{
		MidtransTransactionType:   schedulerReq.TransactionType,
		MidtransTransactionTime:   schedulerReq.TransactionTime,
		MidtransTransactionStatus: schedulerReq.TransactionStatus,
		MidtransTransactionID:     schedulerReq.TransactionID,
		StatusMessage:             schedulerReq.StatusMessage,
		StatusCode:                schedulerReq.StatusCode,
		SignatureKey:              schedulerReq.SignatureKey,
		SettlementTime:            schedulerReq.SettlementTime,
		// ReferenceID:               schedulerReq.ReferenceID,
		PaymentType: schedulerReq.PaymentType,
		OrderID:     schedulerReq.OrderID,
		// Metadata:                  schedulerReq.Metadata,
		MerchantID:  schedulerReq.MerchantID,
		GrossAmount: schedulerReq.GrossAmount,
		FraudStatus: schedulerReq.FraudStatus,
		ExpiryTime:  schedulerReq.ExpiryTime,
		Currency:    schedulerReq.Currency,
		Acquirer:    schedulerReq.Acquirer,
		Body:        body,
	}
}
