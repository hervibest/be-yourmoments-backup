package transactionconsumer

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	errorcode "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/event"
	"github.com/nats-io/nats.go"
)

func (s *TransactionConsumer) setupConsumer(subject string) error {
	consumerConfig := &nats.ConsumerConfig{
		Durable:       s.durableNames[subject],
		AckPolicy:     nats.AckExplicitPolicy,
		MaxDeliver:    5,
		BackOff:       []time.Duration{1 * time.Second, 5 * time.Second, 10 * time.Second},
		DeliverPolicy: nats.DeliverAllPolicy,
		AckWait:       30 * time.Second,
		FilterSubject: subject,
	}

	_, err := s.js.AddConsumer("TRANSACTION_STREAM", consumerConfig)
	return err
}

func (s *TransactionConsumer) handleMessage(ctx context.Context, msg *nats.Msg) {

	var err error
	switch msg.Subject {
	case "transaction.settled":
		s.logs.Log(fmt.Sprintf("received message on subject: %s", msg.Subject))
		event := new(event.OwnerOwnPhotosEvent)
		if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
			_ = msg.Nak()
			s.logs.Error(fmt.Sprintf("failed to unmarshal message : %s", err))
			return
		}

		s.logs.Log(fmt.Sprintf("unmarshalled event: %+v", event))
		request := &model.OwnerOwnPhotosRequest{
			OwnerId:  event.UserId,
			PhotoIds: event.PhotoIds,
		}
		err = s.checkoutUC.OwnerOwnPhotos(ctx, request)
		if err != nil {
			s.handleError(msg, err)
			return
		}

	case "transaction.canceled":
		event := new(event.CancelPhotosEvent)
		if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
			_ = msg.Nak()
			s.logs.Error(fmt.Sprintf("failed to unmarshal message : %s", err))
			return
		}
		request := &model.CancelPhotosRequest{
			UserId:   event.UserId,
			PhotoIds: event.PhotoIds,
		}

		err = s.checkoutUC.CancelPhotos(ctx, request)
		if err != nil {
			s.handleError(msg, err)
			return
		}
	}

	if err := msg.Ack(); err != nil {
		s.logs.Error(fmt.Sprintf("failed to ACK message : %s", err))
	}
}

func (s *TransactionConsumer) handleError(msg *nats.Msg, err error) {
	s.logs.Error(fmt.Sprintf("failed to process transaction event with error : %v ", err))

	appErr, ok := err.(*helper.AppError)
	if !ok {
		appErr = &helper.AppError{Code: errorcode.ErrInternal}
	}

	switch appErr.Code {
	case errorcode.ErrInvalidArgument:
		s.logs.Error(fmt.Sprintf("invalid to ACK message : %s", err))
	default:
		delay := 10 * time.Second
		if err := msg.NakWithDelay(delay); err != nil {
		}
	}
}
