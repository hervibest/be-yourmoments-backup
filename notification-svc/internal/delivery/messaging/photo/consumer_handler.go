package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper"
	errorcode "github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper/enum/error"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/model/event"
	"github.com/nats-io/nats.go"
)

func (s *PhotoConsumer) setupConsumer(subject string) error {
	consumerConfig := &nats.ConsumerConfig{
		Durable:       s.durableNames[subject],
		AckPolicy:     nats.AckExplicitPolicy,
		MaxDeliver:    5,
		BackOff:       []time.Duration{1 * time.Second, 5 * time.Second, 10 * time.Second},
		DeliverPolicy: nats.DeliverAllPolicy,
		AckWait:       30 * time.Second,
		FilterSubject: subject,
	}

	_, err := s.js.AddConsumer("USER_STREAM", consumerConfig)
	return err
}

func (s *PhotoConsumer) handleMessage(ctx context.Context, msg *nats.Msg) {

	var err error
	switch msg.Subject {
	case "photo.bulk":
		event := new(event.BulkPhotoEvent)
		if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
			_ = msg.Nak()
			s.logs.Error(fmt.Sprintf("failed to unmarshal message : %s", err))
			return
		}

		err = s.notficiationUseCase.ProcessAndSendBulkNotificationsV2(ctx, event.UserCountMap)
		if err != nil {
			s.handleError(msg, err, event.EventID)
			return
		}

	case "photo.single.facecam":
		event := new(event.SingleFacecamEvent)
		if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
			_ = msg.Nak()
			s.logs.Error(fmt.Sprintf("failed to unmarshal message : %s", err))
			return
		}

		err = s.notficiationUseCase.ProcessAndSendSingleFacecamNotifications(ctx, event.UserID, event.CountPhotos)
		if err != nil {
			s.handleError(msg, err, event.EventID)
			return
		}

	case "photo.single.photo":
		event := new(event.SinglePhotoEvent)
		if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
			_ = msg.Nak()
			s.logs.Error(fmt.Sprintf("failed to unmarshal message : %s", err))
			return
		}

		err = s.notficiationUseCase.ProcessAndSendSingleNotifications(ctx, event.UserIDs)
		if err != nil {
			s.handleError(msg, err, event.EventID)
			return
		}

	default:
		err = fmt.Errorf("unknown subject: %s", msg.Subject)
	}

	if err := msg.Ack(); err != nil {
		s.logs.Error(fmt.Sprintf("failed to ACK message : %s", err))
	}
}

func (s *PhotoConsumer) handleError(msg *nats.Msg, err error, eventID string) {
	s.logs.Error(fmt.Sprintf("failed to process notification with event id : %s with error : %v ", eventID, err))

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
